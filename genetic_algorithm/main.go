package main

import(
    "math/rand"
)
/* Parameter Restrictions
   - mutation rate in [0,1]
*/


func main() {
    rand.Seed(9)
    lower := 50
    upper := 500
    size := 500
    maxIterations := 500
    lastGen := RunSimulation(size,lower,upper,maxIterations)
    lastGen.WriteToFasta("test.fna")
}
