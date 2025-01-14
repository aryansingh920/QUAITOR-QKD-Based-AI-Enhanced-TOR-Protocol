package middleware

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)


func buildOptimizedRoute(currentPort, finalPort string, trafficData map[string]float64) []string {
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

    // Sort available ports based on traffic data (ascending delay)
    sort.Slice(available, func(i, j int) bool {
        return trafficData[fmt.Sprintf("%d", available[i])] < trafficData[fmt.Sprintf("%d", available[j])]
    })

    // Pick random number of hops: 1–3
    rand.Seed(time.Now().UnixNano())
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






// app.Post("/traffic", func(c *fiber.Ctx) error {
//     var trafficData map[string]float64
//     if err := c.BodyParser(&trafficData); err != nil {
//         return c.Status(fiber.StatusBadRequest).SendString("Invalid data")
//     }

//     // Store traffic data (e.g., in Redis or in-memory map)
//     UpdateTrafficData(trafficData)
//     return c.SendString("Traffic data updated")
// })
