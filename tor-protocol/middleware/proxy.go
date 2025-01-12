package middleware

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// Example: your known network ports
var knownPorts = []int{8801, 8802, 8803, 8804, 8805, 8806, 8807, 8808, 8809, 8810}

// ProxyMiddleware handles forwarding requests to the appropriate server based on the `.onion` path.
// It also supports multi-hop by embedding a route in the "X-Tor-Route" header.
func ProxyMiddleware(c *fiber.Ctx) error {
    // This is the port after /:port.onion/
    // If someone calls http://127.0.0.1:9001/9002.onion/foo
    // then c.Params("port") == "9002"
    // c.Params("*") == "foo"
    finalPort := c.Params("port")
    pathAfterOnion := c.Params("*")

    // Check if there's already a route in the header
    existingRoute := c.Get("X-Tor-Route")

    var route []string
    if existingRoute == "" {
        // No route yet - first node in the chain
        // 1) Build a random route that ends in `finalPort`
        hops := buildRandomRoute(finalPort)
        // 2) Convert slice to a comma-separated string
        route = hops
        log.Printf("[ProxyMiddleware] Generated new route: %v\n", route)
    } else {
        // There's already a route. We'll parse it into a slice.
        route = strings.Split(existingRoute, ",")
    }

    // If for some reason there's nothing in the route (should not happen if built correctly),
    // fallback to the finalPort directly
    if len(route) == 0 {
        route = []string{finalPort}
    }

    // The next hop is the first element in the route
    nextHop := route[0]
    // Remove this hop from the route
    route = route[1:]

    // The updated route (minus the hop we just used)
    updatedRoute := strings.Join(route, ",")

    // If nextHop is the final node, we just forward to that. Otherwise, we keep chaining.
    target := fmt.Sprintf("http://127.0.0.1:%s/%s", nextHop, pathAfterOnion)

    // Set the new route in the request header so the next node sees it
    c.Request().Header.Set("X-Tor-Route", updatedRoute)

    // Forward to next hop
    return proxy.Forward(target)(c)
}

// ProxyExactMiddleware handles requests with no trailing slash for `.onion` routes.
// You can similarly embed the multi-hop logic here if you want the same chain approach.
func ProxyExactMiddleware(c *fiber.Ctx) error {
    finalPort := c.Params("port")

    // For the no-path scenario, we do the same route building or usage as above.
    // In many cases, you might unify these into one function, but for clarity, keep them separate.
    existingRoute := c.Get("X-Tor-Route")

    var route []string
    if existingRoute == "" {
        route = buildRandomRoute(finalPort)
        log.Printf("[ProxyExactMiddleware] Generated new route: %v\n", route)
    } else {
        route = strings.Split(existingRoute, ",")
    }

    if len(route) == 0 {
        route = []string{finalPort}
    }

    nextHop := route[0]
    route = route[1:]
    updatedRoute := strings.Join(route, ",")

    target := fmt.Sprintf("http://127.0.0.1:%s", nextHop)
    c.Request().Header.Set("X-Tor-Route", updatedRoute)

    return proxy.Forward(target)(c)
}

// buildRandomRoute constructs a random route of length, for example, 3 or 4 intermediate hops,
// eventually ending in finalPort. Tweak as needed.
func buildRandomRoute(finalPort string) []string {
    // E.g., pick 2 or 3 random nodes from knownPorts excluding the final port (and possibly the current port).
    // Then append finalPort at the end.

    // Convert finalPort from string to int just for safety
    fPort, _ := strconv.Atoi(finalPort)

    // We don't want to include the final port among the random “intermediate” hops
    var available []int
    for _, p := range knownPorts {
        if p != fPort {
            available = append(available, p)
        }
    }

    // Shuffle your available ports
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(available), func(i, j int) {
        available[i], available[j] = available[j], available[i]
    })

    // Pick some random number of hops, e.g. 2 or 3
    numHops := rand.Intn(3) + 1 // random 1–3 hops
    if numHops > len(available) {
        numHops = len(available)
    }
    intermediateHops := available[:numHops]

    // Build the route
    route := make([]string, 0, numHops+1)
    // Convert each intermediate hop to string
    for _, hop := range intermediateHops {
        route = append(route, fmt.Sprintf("%d", hop))
    }

    // Finally, append the final port
    route = append(route, finalPort)

    return route
}
