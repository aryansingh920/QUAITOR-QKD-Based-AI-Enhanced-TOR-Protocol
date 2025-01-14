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
        route = buildRandomRoute(currentPort, finalPort)
        log.Printf("[Port %s] [ProxyMiddleware] Generated new route: %v\n", currentPort, route)
    } else {
        route = strings.Split(existingRoute, ",")
        log.Printf("[Port %s] [ProxyMiddleware] Existing route: %v\n", currentPort, route)
    }

    // If no route left, target is the finalPort
    if len(route) == 0 {
        route = []string{finalPort}
    }

    nextHop := route[0]
    route = route[1:]
    updatedRoute := strings.Join(route, ",")

    // --- 1) Parse and potentially DECRYPT incoming query params ---
    incomingQuery := c.Request().URI().QueryString()
    queryParams := parseQueryParams(string(incomingQuery))

    // Let's assume only "msg" and "entry" might be encrypted
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

    // --- 2) (Optional) Do something with the decrypted parameters here ---

    // --- 3) Potentially ENCRYPT before sending to next hop ---
    // If nextHop is NOT the final hop, you might want to re-encrypt.
    // If nextHop == finalPort (and it's truly final), you might send plaintext.
    // This example encrypts again unconditionally (for demonstration).
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

    // Now rebuild the query string
    newQueryString := buildQueryString(queryParams)

    // Build the target URL
    target := fmt.Sprintf("%s:%s/%s", config.DefaultLink, nextHop, pathAfterOnion)
    if newQueryString != "" {
        target += "?" + newQueryString
    }

    // Simulate a random delay
    randomDelay := time.Duration(rand.Intn(config.RandomDelayUpperLimit)) * time.Millisecond
    log.Printf("[Port %s] Adding random delay: %s", currentPort, randomDelay)
    time.Sleep(randomDelay)

    log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
        currentPort, c.IP(), nextHop, route)

    // Update the route header
    c.Request().Header.Set(config.CustomHeaderKey, updatedRoute)

    // Forward the request
    return proxy.Forward(target)(c)
}


// parseQueryParams converts a raw query string into a map of key/value pairs.
// This is a simplistic implementation; you can also use url.ParseQuery for more robust handling.
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



// func ProxyMiddleware(c *fiber.Ctx) error {
//     currentPort := config.GetPort()
//     finalPort := c.Params("port")
//     pathAfterOnion := c.Params("*")
//     existingRoute := c.Get(config.CustomHeaderKey)
//     var route []string

//     if existingRoute == "" {
//         route = buildRandomRoute(currentPort, finalPort)
//         log.Printf("[Port %s] [ProxyMiddleware] Generated new route: %v\n", currentPort, route)
//     } else {
//         route = strings.Split(existingRoute, ",")
//         log.Printf("[Port %s] [ProxyMiddleware] Existing route: %v\n", currentPort, route)
//     }

//     if len(route) == 0 {
//         route = []string{finalPort}
//     }

//     nextHop := route[0]
//     route = route[1:]
//     updatedRoute := strings.Join(route, ",")

//     // Build the base target
//     target := fmt.Sprintf("%s:%s/%s", config.DefaultLink, nextHop, pathAfterOnion)

//     // Append the query string if it exists
//     qs := c.Request().URI().QueryString()
//     if len(qs) > 0 {
//         target = target + "?" + string(qs)
//     }

//     // Simulate processing delay at the current node
//     randomDelay := time.Duration(rand.Intn(config.RandomDelayUpperLimit)) * time.Millisecond
//     log.Printf("[Port %s] Adding random delay: %s", currentPort, randomDelay)
//     time.Sleep(randomDelay)

//     log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
//         currentPort, c.IP(), nextHop, route)

//     // Update the route header
//     c.Request().Header.Set(config.CustomHeaderKey, updatedRoute)

//     // Now forward the request (including the original query string)
//     return proxy.Forward(target)(c)
// }

// func ProxyMiddleware(c *fiber.Ctx) error {
//     currentPort := config.GetPort()
//     finalPort := c.Params("port")
//     pathAfterOnion := c.Params("*")
//     existingRoute := c.Get(config.CustomHeaderKey)
//     var route []string

//     if existingRoute == "" {
//         route = buildRandomRoute(currentPort, finalPort)
//         log.Printf("[Port %s] [ProxyMiddleware] Generated new route: %v\n", currentPort, route)
//     } else {
//         route = strings.Split(existingRoute, ",")
//         log.Printf("[Port %s] [ProxyMiddleware] Existing route: %v\n", currentPort, route)
//     }

//     if len(route) == 0 {
//         route = []string{finalPort}
//     }

//     nextHop := route[0]
//     route = route[1:]
//     updatedRoute := strings.Join(route, ",")

//     target := fmt.Sprintf("%s:%s/%s",config.DefaultLink, nextHop, pathAfterOnion)

//     // Simulate processing delay at the current node
//     randomDelay := time.Duration(rand.Intn(config.RandomDelayUpperLimit)) * time.Millisecond // Up to 5 seconds
//     log.Printf("[Port %s] Adding random delay: %s", currentPort, randomDelay)
//     time.Sleep(randomDelay)

//     log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
//         currentPort, c.IP(), nextHop, route)

//     c.Request().Header.Set(config.CustomHeaderKey, updatedRoute)
//     return proxy.Forward(target)(c)
// }

// ProxyExactMiddleware handles requests with no trailing slash for .onion routes.
// func ProxyExactMiddleware(c *fiber.Ctx) error {
// 	currentPort := config.GetPort()
// 	finalPort := c.Params("port")

// 	existingRoute := c.Get(config.CustomHeaderKey)
// 	var route []string

// 	if existingRoute == "" {
// 		route = buildRandomRoute(currentPort, finalPort)
// 		log.Printf("[Port %s] [ProxyExactMiddleware] Generated new route: %v\n", currentPort, route)
// 	} else {
// 		route = strings.Split(existingRoute, ",")
// 		log.Printf("[Port %s] [ProxyExactMiddleware] Existing route: %v\n", currentPort, route)
// 	}

// 	if len(route) == 0 {
// 		route = []string{finalPort}
// 	}

// 	nextHop := route[0]
// 	route = route[1:]
// 	updatedRoute := strings.Join(route, ",")

// 	target := fmt.Sprintf("%s:%s",config.DefaultLink, nextHop)

// 	log.Printf("[Port %s] Received request from: %s => next hop: %s => remaining route: %v\n",
// 		currentPort, c.IP(), nextHop, route)

// 	c.Request().Header.Set(config.CustomHeaderKey, updatedRoute)
// 	return proxy.Forward(target)(c)
// }

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
