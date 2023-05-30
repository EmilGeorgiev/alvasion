# alvasion
This is a simulation of an alien invasion.

### Additional Assumptions To The Task
- A city cannot be destroyed by a lone alien; a minimum of two aliens must concurrently inhabit a city to cause its destruction. 
- If more than two aliens simultaneously inhabit a city, the city will be demolished and the aliens will be annihilated.
- Initially, when aliens are dispersed randomly among the cities, it is guaranteed that no city will accommodate more than one alien. 
- If the number of aliens surpasses the number of cities, aliens will be distributed until each city hosts a single alien; any remaining aliens will not be assigned to any city.
- The city connections are accurately outlined in the provided file. For instance, if a connection is specified as "X2 west=X1", you can be certain that there will also be a connection stated as "X1 east=X2"."
- The next iteration doesn't start until all aliens finished with their movements and the commander made has updates on the map 

## Invasion Workflow
The invasion contains several steps: "Generate The World Map", "Create Alien Commander", "Start the invasion" and "Commander write a report when invasion finished"

### Generate The World Map
This segment of the program initializes the invasion process by constructing a world map based on data from a specified 
file containing city information and their interconnected routes. To accomplish this task effectively, the program 
leverages concurrent programming with several goroutines.

Initially, a single goroutine reads the file line by line, forwarding each line into a dedicated channel. This task follows 
the Fan-Out pattern, where one goroutine disseminates data through a channel to multiple goroutines. The number of goroutines 
reading from this channel, called 'workers', can be configured based on the requirements. Each worker is responsible for 
validating and splitting the incoming line into distinct parts using a whitespace as the delimiter.

Once the lines have been validated, the workers send these lines to another goroutine responsible for parsing and converting 
them into a format suitable for world map creation. This phase utilizes the Fan-In pattern, characterized by multiple 
goroutines feeding data through a channel into a single goroutine.

This design approach is beneficial when the source file contains a large number of lines. By leveraging concurrent processing, 
the program can efficiently parse, validate, and transform data, providing an optimized way to generate the world map for 
the invasion process.
 
```
                                   |---------------|   
                                   | validate line | 
                                   |---------------|
          
 |-----------------------|   Line  |---------------|  parts of the line    |--------------|
 |read file line by line | ------> | validate line | ------------------->  | generate map |
 |-----------------------|         |---------------|                       |--------------|

                                   |---------------|   
                                   | validate line | 
                                   |---------------|

```

### Create Alien Commander
AlienCommander serves as the strategic leader and coordinator of the alien forces during an invasion.
This entity manages the distribution and movement of alien soldiers across the global map, delegating orders
on which outgoing roads soldiers should advance through in a city.

Equipped with a comprehensive world map and a roster of all alien soldiers at its disposal, the AlienCommander
maintains real-time awareness of the location and status of each soldier. This includes tracking soldiers that
have been killed or trapped during the course of the invasion.

Besides managing the alien troops, the AlienCommander is also in charge of controlling the overall progression
of the invasion. It triggers each iteration of the invasion cycle, making decisive calls on when to start the next
invasion iteration or when to cease operations entirely.

Additionally, the AlienCommander oversees the evaluation of each city post-invasion and responds to situation
reports (Sitreps) from soldiers, allowing it to promptly react to the changing dynamics of the invasion landscape.

In essence, the AlienCommander plays a pivotal role as the central management unit of the invasion, responsible
for executing strategic decisions, commanding soldiers, and maintaining the pace of the invasion.

### Start The Invasion
The invasion is started by AlienCommander. The invasion is separated on several steps:
The invasion process is initiated by the AlienCommander and unfolds over several distinct steps:

1. Alien Distribution: In the initial stage, the commander strategically disperses his alien soldiers across various cities, with a limitation of one alien per city. This operation is performed only once at the beginning of the invasion. All subsequent operations occur in every iteration of the invasion process.
2. Command Issuance: The commander dictates the direction of each soldier's advance, choosing from North, South, East, or West. The selection of direction is determined randomly, creating a unique and unpredictable pattern of invasion.
3. Invasion Iteration: Upon execution of the given orders, the commander begins evaluating the impact of the invasion. During this phase, the commander actively monitors reports from the soldiers, updates the world map accordingly, and oversees the strategic demolition of cities based on the progress of the invasion.
4. Evaluation of Invasion Status: At this stage, the commander assesses the ongoing state of the invasion, determining whether it should be terminated or continued. The invasion is typically concluded when there's either zero or one alien soldier left, or when each soldier has completed 10,000 iterations of the invasion process.
5. Iterative Invasion Progression: Following the evaluation, the process circles back to the 'Command Issuance' stage. Steps 2, 3, and 4 are repeated in iterative cycles, thereby advancing the invasion until the predetermined conclusion criteria are met.

With this cyclic and strategic approach, the AlienCommander manages the invasion, dynamically adjusting to evolving conditions, and making informed decisions that drive the alien forces towards their objective.

### Write a report when invasion finished
This stage represents the culmination of the program's execution. Once the invasion has concluded, the AlienCommander compiles a comprehensive report detailing the aftermath. The report outlines the remaining cities and their interconnected roads, providing a clear view of the post-invasion landscape. This detailed document, serving as the official record of the invasion's outcome, is subsequently stored in a dedicated file.

## Run the program
You can run the program locally.

### Run the program locally
First clone the repo:
```
git clone git@github.com:EmilGeorgiev/alvasion.git
```

Navigate to the "alvasion/cmd" folder:
```
cd alvasion/cmd
```

For testing purposes we have a file world-map.txt in the cmd fodler that contains all cities in the world and their connections
the file contains these cities (the file already exist and is not necessary to crated it, but if you want you can provide your own data in it): 
```
X1 east=X2 south=X4
X2 east=X3 west=X1 south=X5
X3 west=X2 south=X6
X4 east=X5 north=X1 south=X7
X5 west=X4 east=X6 north=X2 south=X8
X6 west=X5 north=X3 south=X9
X7 east=X8 north=X4
X8 west=X7 east=X9 north=X5
X9 west=X8 north=X6
```

Run the project:
```
go run main.go -aliens=6
```

You should see result similar to:
```
2023/05/30 10:57:16 Generating World Map.
2023/05/30 10:57:16 WorldMap is generated.
2023/05/30 10:57:16 Initialize 6 number of aliens/soldiers.
2023/05/30 10:57:16 Initialize AlienCommander.
2023/05/30 10:57:16 Start the invasion!
X5 is destroyed from alien 0 and alien 1!
X8 is destroyed from alien 3 and alien 2!
X2 is destroyed from alien 5 and alien 4!
There is 0 soldier left. Stop the invasion!
Number of iterations:  8
2023/05/30 10:57:16 Generate the report
2023/05/30 10:57:16 Store the report in a file report.txt
2023/05/30 10:57:16 Finish
```

In the file aliens/cmd/report.txt - you can see the report of the invasion. It contains information something like this:
```
X1 south=X4
X3 south=X6
X4 north=X1 south=X7
X6 north=X3 south=X9
X7 north=X4
X9 north=X6
```

### Configurations
The project contains a configuration file located in: ./cmd/config.yaml. In this file you can configure
where the world-map.txt file is and the number of validation workers that will validate the lines of the file.
Later here can be added a new configuration variables.