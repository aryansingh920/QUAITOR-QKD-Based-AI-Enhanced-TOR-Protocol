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

// Define port range for knownPorts generation
const (
	portStart = 8801
	portEnd   = 8820
)

// getKnownPorts dynamically generates the list of known ports based on the range.
func getKnownPorts() []int {
	var ports []int
	for p := portStart; p <= portEnd; p++ {
		ports = append(ports, p)
	}
	return ports
}

// ProxyMiddleware handles forwarding requests to the appropriate server based on the .onion path
// or based on an existing X-Tor-Route header. It supports multi-hop by embedding a route in
// the "X-Tor-Route" header.
func ProxyMiddleware(c *fiber.Ctx) error {
	// Current node's port (the one we are *running* on)
	currentPort := config.GetPort()

	// The requested final port from URL pattern /:port.onion/
	// or /:port/ if user visits /8805/ or /8805.onion
	finalPort := c.Params("port")

	// The path after the .onion e.g. "/foo"
	pathAfterOnion := c.Params("*")

	// Retrieve the existing route from the header (if any)
	existingRoute := c.Get("X-Tor-Route")

	var route []string

	// 1) If there's no existingRoute, we create a brand new route that eventually ends in finalPort
	if existingRoute == "" {
		// Means this is the first node in the chain
		route = buildRandomRoute(currentPort, finalPort)
		log.Printf("[Port %s] [ProxyMiddleware] Generated new route: %v\n", currentPort, route)
	} else {
		// There's already a route
		route = strings.Split(existingRoute, ",")
		log.Printf("[Port %s] [ProxyMiddleware] Existing route: %v\n", currentPort, route)
	}

	// Safety check
	if len(route) == 0 {
		// If for some reason route is empty, fallback to finalPort only
		route = []string{finalPort}
	}

	// The next hop is always route[0]
	nextHop := route[0]

	log.Printf("[Port %s] Current X-Tor-Route: %s", currentPort, existingRoute)
	log.Printf("[Port %s] Remaining Route: %v", currentPort, route)
	log.Printf("[Port %s] Next Hop: %s", currentPort, nextHop)

	// Remove this hop from the route
	route = route[1:]

	// The updated route for the next node
	updatedRoute := strings.Join(route, ",")

	// Build the forwarding URL
	//   e.g. "http://127.0.0.1:8810/foo"
	//   pathAfterOnion might be "" if the user visited "/8805.onion/"
	//   so effectively forwarding to 8810/
	target := fmt.Sprintf("http://127.0.0.1:%s/%s", nextHop, pathAfterOnion)

	// Log the hop
	log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
		currentPort, c.IP(), nextHop, route)

	// Set the new route in the request header so the next node sees it
	c.Request().Header.Set("X-Tor-Route", updatedRoute)

	log.Printf("[Port %s] Updated X-Tor-Route: %s", currentPort, updatedRoute)

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
		route = buildRandomRoute(currentPort, finalPort)
		log.Printf("[Port %s] [ProxyExactMiddleware] Generated new route: %v\n", currentPort, route)
	} else {
		route = strings.Split(existingRoute, ",")
		log.Printf("[Port %s] [ProxyExactMiddleware] Existing route: %v\n", currentPort, route)
	}

	if len(route) == 0 {
		route = []string{finalPort}
	}

	nextHop := route[0]
	route = route[1:]
	updatedRoute := strings.Join(route, ",")

	target := fmt.Sprintf("http://127.0.0.1:%s", nextHop)

	log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
		currentPort, c.IP(), nextHop, route)

	c.Request().Header.Set("X-Tor-Route", updatedRoute)
	return proxy.Forward(target)(c)
}

// buildRandomRoute constructs a random route (1–3 intermediate hops), excluding the currentPort
// and excluding finalPort as an intermediate hop. The final hop is always finalPort.
func buildRandomRoute(currentPort, finalPort string) []string {
	fPort, _ := strconv.Atoi(finalPort)
	cPort, _ := strconv.Atoi(currentPort)

	// Generate the known ports list
	allPorts := getKnownPorts()

	// Filter out currentPort and finalPort from available ports
	var available []int
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

	// Select intermediate hops
	intermediateHops := available[:numHops]

	// Build the route array (all intermediate hops + final port)
	route := make([]string, 0, numHops+1)
	for _, hop := range intermediateHops {
		route = append(route, fmt.Sprintf("%d", hop))
	}
	route = append(route, finalPort)

	return route
}
