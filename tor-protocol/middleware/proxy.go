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
var (
	portStart = config.PortStart
	portEnd   = config.PortEnd
)

// getKnownPorts dynamically generates the list of known ports based on the range.

// parseQueryParams converts a raw query string into a map of key/value pairs.

// buildQueryString rebuilds a query string from a map.

func ProxyExactMiddleware(c *fiber.Ctx) error {
	currentPort := config.GetPort()
	finalPort := c.Params("port")

	existingRoute := c.Get(config.CustomHeaderKey)
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

	// Build the base target
	target := fmt.Sprintf("%s:%s", config.DefaultLink, nextHop)

	// Append the query string if it exists
	qs := c.Request().URI().QueryString()
	if len(qs) > 0 {
		target = target + "?" + string(qs)
	}

	log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
		currentPort, c.IP(), nextHop, route)

	// Update the route header
	c.Request().Header.Set(config.CustomHeaderKey, updatedRoute)

	err := proxy.Forward(target)(c)
	if err != nil {
		log.Printf("[Port %s] Error forwarding to %s: %v", currentPort, target, err)
		return err
	}
	return nil

	// return proxy.Forward(target)(c)
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
	fmt.Printf("Available ports: %v\n", available)
	fmt.Printf("Port start: %d, Port end: %d\n", portStart, portEnd)

	genHops := portEnd - portStart
	numHops := rand.Intn(genHops) + 1
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

// getKnownPorts dynamically generates the list of known ports based on the range.
func getKnownPorts() []int {
	var ports []int
	for p := portStart; p <= portEnd; p++ {
		ports = append(ports, p)
	}
	return ports
}

func ProxyMiddleware(c *fiber.Ctx) error {
    currentPort := config.GetPort()
    finalPort := c.Params("port")
    pathAfterOnion := c.Params("*")
    existingRoute := c.Get(config.CustomHeaderKey)
    var route []string

    // Build or retrieve existing route
    if existingRoute == "" {
        // We are the *first node* in the chain
        route = buildRandomRoute(currentPort, finalPort)
        log.Printf("[Port %s] [ProxyMiddleware] Generated new route: %v\n", currentPort, route)
    } else {
        // We are an *intermediate* or *last* node
        route = strings.Split(existingRoute, ",")
        log.Printf("[Port %s] [ProxyMiddleware] Existing route: %v\n", currentPort, route)
    }

    // --------------------------------------------------------------------
    // 1) If the route is empty *already*, we are truly last → decrypt + no forward
    // --------------------------------------------------------------------
    if len(route) == 0 {
        // This node is *the final node* (no nextHop).
        log.Printf("[Port %s] This is the FINAL node (route empty). Decrypting + local handling.\n", currentPort)

        // Decrypt "msg" and "entry" if they exist
        incomingQuery := c.Request().URI().QueryString()
        queryParams := parseQueryParams(string(incomingQuery))

        if val, ok := queryParams["msg"]; ok && val != "" {
            decrypted, err := decryptMessage(val)
            if err != nil {
                log.Printf("[Port %s] Decryption error for 'msg': %v", currentPort, err)
            } else {
                queryParams["msg"] = decrypted
            }
        }
        if val, ok := queryParams["entry"]; ok && val != "" {
            decrypted, err := decryptMessage(val)
            if err != nil {
                log.Printf("[Port %s] Decryption error for 'entry': %v", currentPort, err)
            } else {
                queryParams["entry"] = decrypted
            }
        }

        // Optionally, rewrite the Fiber query so downstream handlers see the decrypted text:
        newQueryString := buildQueryString(queryParams)
        c.Request().URI().SetQueryString(newQueryString)

        // Return c.Next() so your local final handler (e.g. controllers.HomeHandler) can serve it:
        return c.Next()
    }

    // --------------------------------------------------------------------
    // 2) Otherwise, we are not last. We pop the next hop + check if we’re first node for encryption
    // --------------------------------------------------------------------
    nextHop := route[0]
    remaining := route[1:] // route after popping nextHop
    updatedRoute := strings.Join(remaining, ",")

    // parse incoming query
    incomingQuery := c.Request().URI().QueryString()
    queryParams := parseQueryParams(string(incomingQuery))

    // Are we the first node in the entire chain (existingRoute == "")?
    isFirstNode := (existingRoute == "")

    // If first node, ENCRYPT
    if isFirstNode {
        if val, ok := queryParams["msg"]; ok && val != "" {
            encrypted, err := encryptMessage(val)
            if err != nil {
                log.Printf("[Port %s] Encryption error for 'msg': %v", currentPort, err)
            } else {
                queryParams["msg"] = encrypted
            }
        }
        if val, ok := queryParams["entry"]; ok && val != "" {
            encrypted, err := encryptMessage(val)
            if err != nil {
                log.Printf("[Port %s] Encryption error for 'entry': %v", currentPort, err)
            } else {
                queryParams["entry"] = encrypted
            }
        }
    }

    // Rebuild the query string after possible encryption
    newQueryString := buildQueryString(queryParams)

    // Build the target URL for nextHop
    target := fmt.Sprintf("%s:%s/%s", config.DefaultLink, nextHop, pathAfterOnion)
    if newQueryString != "" {
        target += "?" + newQueryString
    }

    // Simulate a random delay
    randomDelay := time.Duration(rand.Intn(config.RandomDelayUpperLimit)) * time.Millisecond
    log.Printf("[Port %s] Adding random delay: %s", currentPort, randomDelay)
    time.Sleep(randomDelay)

    log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
        currentPort, c.IP(), nextHop, remaining)

    // Update the X-Tor-Route header so next node sees the updated route
    c.Request().Header.Set(config.CustomHeaderKey, updatedRoute)

    // Forward the request
    if err := proxy.Forward(target)(c); err != nil {
        log.Printf("[Port %s] Error forwarding to %s: %v", currentPort, target, err)
        return err
    }
    return nil
}


// parseQueryParams converts a raw query string into a map of key/value pairs.
func parseQueryParams(rawQuery string) map[string]string {
	result := make(map[string]string)
	if rawQuery == "" {
		return result
	}
	pairs := strings.Split(rawQuery, "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := parts[1]
			result[key] = val
		} else {
			// Handle keys with no value if needed
			result[parts[0]] = ""
		}
	}
	return result
}

// buildQueryString rebuilds a query string from a map.
func buildQueryString(params map[string]string) string {
	var parts []string
	for k, v := range params {
		// URL-encode them if needed
		part := fmt.Sprintf("%s=%s", k, v)
		parts = append(parts, part)
	}
	return strings.Join(parts, "&")
}
