// middleware.go
package middleware

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"tor-protocol/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// Known network ports
// You can adapt this to whatever known onion-like nodes (ports) you have.
// Define port range for knownPorts generation
const (
    portStart = 8801
    portEnd   = 8810
)

// getKnownPorts dynamically generates the list of known ports based on the range.
func getKnownPorts() []int {
    var ports []int
    for p := portStart; p <= portEnd; p++ {
        ports = append(ports, p)
    }
    return ports
}

// ProxyMiddleware handles forwarding requests to the appropriate server based on the .onion path.
// It supports multi-hop by embedding a route in the "X-Tor-Route" header.
func ProxyMiddleware(c *fiber.Ctx) error {
    // Current node's port (the one we are *running* on)
    // e.g., 8807 if we started with `go run main.go 8807`
    currentPort := config.GetPort()

    // The requested final port from URL pattern /:port.onion/
    finalPort := c.Params("port")

    // The path after the .onion e.g. "/foo"
    pathAfterOnion := c.Params("*")

    // Retrieve the existing route from the header
    existingRoute := c.Get("X-Tor-Route")

    var route []string

    if existingRoute == "" {
        // No route yet => first node in the chain
        // Build a random route that ends in finalPort, excluding currentPort
        route = buildRandomRoute(currentPort, finalPort)
        log.Printf("[Port %s] [ProxyMiddleware] Generated new route: %v\n", currentPort, route)
    } else {
        // There's already a route
        route = strings.Split(existingRoute, ",")
        log.Printf("[Port %s] [ProxyMiddleware] Existing route: %v\n", currentPort, route)
    }

    // If for some reason there's nothing in the route (should not happen if built correctly),
    // fallback to just finalPort
    if len(route) == 0 {
        route = []string{finalPort}
    }

    // The next hop is always route[0]
    nextHop := route[0]
    // Remove this hop from the route
    route = route[1:]

    // The updated route (minus the hop we just used)
    updatedRoute := strings.Join(route, ",")

    // If nextHop is the final node, we just forward to that. Otherwise, we keep chaining.
    target := fmt.Sprintf("http://127.0.0.1:%s/%s", nextHop, pathAfterOnion)

    // Log the hop
    log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
        currentPort, c.IP(), nextHop, route)

    // Set the new route in the request header so the next node sees it
    c.Request().Header.Set("X-Tor-Route", updatedRoute)

    // Forward to next hop
    return proxy.Forward(target)(c)
}

// ProxyExactMiddleware handles requests with no trailing slash for .onion routes.
func ProxyExactMiddleware(c *fiber.Ctx) error {
    currentPort := config.GetPort()
    finalPort := c.Params("port")

    existingRoute := c.Get("X-Tor-Route")
    var route []string

    if existingRoute == "" {
        // Build a new route
        route = buildRandomRoute(currentPort, finalPort)
        log.Printf("[Port %s] [ProxyExactMiddleware] Generated new route: %v\n", currentPort, route)
    } else {
        // There's already a route
        route = strings.Split(existingRoute, ",")
        log.Printf("[Port %s] [ProxyExactMiddleware] Existing route: %v\n", currentPort, route)
    }

    if len(route) == 0 {
        route = []string{finalPort}
    }

    nextHop := route[0]
    route = route[1:]
    updatedRoute := strings.Join(route, ",")

    // If nextHop is the final node, we just forward
    target := fmt.Sprintf("http://127.0.0.1:%s", nextHop)

    log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
        currentPort, c.IP(), nextHop, route)

    c.Request().Header.Set("X-Tor-Route", updatedRoute)
    return proxy.Forward(target)(c)
}

// buildRandomRoute constructs a random route of length 1–3 intermediate hops (you can tweak this)
// excluding the currentPort and excluding finalPort as an intermediate hop.
func buildRandomRoute(currentPort, finalPort string) []string {
    // Convert finalPort from string to int
    fPort, _ := strconv.Atoi(finalPort)
    cPort, _ := strconv.Atoi(currentPort)

    // We don't want the finalPort as an intermediate node,
    // and we don't want to re-use the currentPort as an intermediate node
    var available []int
    allPorts := getKnownPorts()
    for _, p := range allPorts {
        if p != fPort && p != cPort {
            available = append(available, p)
        }
    }

    // Shuffle the available ports
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(available), func(i, j int) {
        available[i], available[j] = available[j], available[i]
    })

    // Pick random number of hops: 1–3
    numHops := rand.Intn(3) + 1
    if numHops > len(available) {
        numHops = len(available)
    }

    // intermediateHops are the random selection
    intermediateHops := available[:numHops]

    // Build the route array
    route := make([]string, 0, numHops+1)
    for _, hop := range intermediateHops {
        route = append(route, fmt.Sprintf("%d", hop))
    }

    // Finally, append the final port as the last hop
    route = append(route, finalPort)

    return route
}
