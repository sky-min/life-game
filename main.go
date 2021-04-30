package main

import (
    "fmt"
    "math/rand"
    "time"
    "os"
    "os/exec"
)

const W = 40
const H = 40

const LIVE = 1
const DEAD = 0

type Field [W][H]int

type Traverser func(cell *int, x, y, maxX, maxY int)

func Traverse(f *Field, fn Traverser) {
    for i := 0; i < len(f); i++ {
        for j := 0; j < len(f[i]); j++ {
            fn(&f[i][j], i, j, len(f) - 1, len(f[i]) - 1)
        }
    }
}

func DrawField(f *Field) {
    // linux only :(
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()

    Traverse(f, func(cell *int, x, y, maxX, maxY int) {
        var c string

        if 0 == *cell {
            c = "."
        } else {
            c = "X"
        }

        fmt.Printf("%s", c)

        if y == maxY {
            fmt.Printf("\n")
        }
    })
}

func FillField(f *Field) {
    rand.Seed(time.Now().UTC().UnixNano())

    Traverse(f, func(cell *int, x, y, maxX, maxY int) {
        if 1 == rand.Intn(4) {
            *cell = LIVE
        }
    })    
}

func FillFieldGlider(f *Field) {
    f[5][3] = LIVE
    f[5][4] = LIVE
    f[5][5] = LIVE
    f[4][5] = LIVE
    f[3][4] = LIVE
}

func Step(gen Field) Field {
    var x, y int
    var nextGen Field

    for i := 0; i < len(gen); i++ {
        for j := 0; j < len(gen[i]); j++ {
            nNear := 0

            for ii := -1; ii < 2; ii++ {
                for jj := -1; jj < 2; jj++ {
                    if 0 == ii && 0 == jj {
                        continue
                    }

                    x, y = i + ii, j + jj

                    if 0 <= x && x < len(gen) && 0 <= y && y < len(gen[i]) && LIVE == gen[x][y] {
                        nNear++
                    }
                }
            }
            
            nextGen[i][j] = MakeCell(nNear, LIVE == gen[i][j]) 
        }
    }

    return nextGen
}

func MakeCell(nNear int, isLive bool) int {
    if isLive && (2 == nNear || 3 == nNear) {
        return LIVE
    } 
    
    if !isLive && 3 == nNear {
        return LIVE
    }

    return DEAD
}

func main() {
    var f Field

    FillFieldGlider(&f)

    for {
        DrawField(&f)
        f = Step(f)
        time.Sleep(300 * time.Millisecond)
    }
}