package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"

    "golang.org/x/exp/slices"
)

const (
    numPhilo       = 5
    eatTimes       = 3
    numEatingPhilo = 2
)

type eatRequest struct {
    who            int         // Who is making the request
    finishedFnChan chan func() // When approves a response will be sent on this channel with a function to call when done
}

// simulateHost - the host must provide permission before a philosopher can eat
// Exits when channel closed
func simulateHost(requestChannel <-chan eatRequest) {
    awaitRequest := requestChannel
    finishedChan := make(chan struct {
        who  int
        done chan struct{}
    })

    var whoEating []int // tracks who is currently eating

    for {
        select {
        case request, ok := <-awaitRequest:
            if !ok {
                return // Closed channel means that we are done (finishedChan is guaranteed to be empty)
            }
            // Sanity check - confirm that philosopher is not being greedy! (should never happen)
            if slices.Index(whoEating, request.who) != -1 {
                panic("Multiple requests from same philosopher")
            }
            whoEating = append(whoEating, request.who) // New request always goes at the end
            fmt.Printf("%d started eating (currently eating %v)\n", request.who, whoEating)

            // Let philosopher know and provide means for them to tell us when done
            request.finishedFnChan <- func() {
                d := make(chan struct{})
                finishedChan <- struct {
                    who  int
                    done chan struct{}
                }{who: request.who, done: d}
                <-d // Wait until request has been processed (ensure we should never have two active requests from one philosopher)
            }
        case fin := <-finishedChan:
            idx := slices.Index(whoEating, fin.who)
            if idx == -1 {
                panic("philosopher stopped eating multiple times!")
            }
            whoEating = append(whoEating[:idx], whoEating[idx+1:]...) // delete the element
            fmt.Printf("%d completed eating (currently eating %v)\n", fin.who, whoEating)
            close(fin.done)
        }
        // There has been a change in the number of philosopher's eating
        if len(whoEating) < numEatingPhilo {
            awaitRequest = requestChannel
        } else {
            awaitRequest = nil // Ignore new eat requests until a philosopher finishes (nil channel will never be selected)
        }
    }
}

// ChopS represents a single chopstick
type ChopS struct {
    mu  sync.Mutex
    idx int // Including the index can make debugging simpler
}

// philosopher simulates a Philosopher (brain in a vat!)
func philosopher(philNum int, leftCS, rightCS *ChopS, requestToEat chan<- eatRequest) {
    for numEat := 0; numEat < eatTimes; numEat++ {
        // once the philosopher intends to eat, lock the corresponding chopsticks
        for {
            leftCS.mu.Lock()
            // Attempt to get the right Chopstick - if someone else has it we replace the left chopstick and try
            // again (in order to avoid deadlocks)
            if rightCS.mu.TryLock() {
                break
            }
            leftCS.mu.Unlock()
        }

        // We have the chopsticks but need the hosts permission
        ffc := make(chan func()) // when accepted we will receive a function to call when done eating
        requestToEat <- eatRequest{
            who:            philNum,
            finishedFnChan: ffc,
        }
        doneEating := <-ffc

        fmt.Printf("philosopher %d starting to eat (%d feed)\n", philNum, numEat)
        time.Sleep(time.Millisecond * time.Duration(rand.Intn(200))) // Eating takes a random amount of time
        fmt.Printf("philosopher %d finished eating (%d feed)\n", philNum, numEat)

        rightCS.mu.Unlock()
        leftCS.mu.Unlock()
        doneEating() // Tell host that we have finished eating
    }
    fmt.Printf("philosopher %d is full\n", philNum)
}

func main() {
    CSticks := make([]*ChopS, numPhilo)
    for i := 0; i < numPhilo; i++ {
        CSticks[i] = &ChopS{idx: i}
    }

    requestChannel := make(chan eatRequest)

    var wg sync.WaitGroup
    wg.Add(numPhilo)
    for i := 0; i < numPhilo; i++ {
        go func(philNo int) {
            philosopher(philNo, CSticks[philNo-1], CSticks[philNo%numPhilo], requestChannel)
            wg.Done()
        }(i + 1)
    }

    go func() {
        wg.Wait()
        close(requestChannel)
    }()

    simulateHost(requestChannel)
}