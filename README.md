# Word of Wisdom TCP server

![Sage](https://www.meme-arsenal.com/memes/e2064a411ac3fa1066f90ae854a64d60.jpg)

- TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.              
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.              
- Docker file should be provided both for the server and for the client that solves the POW challenge.

## Algorithm
[Guided Tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol) (Network-bound) algorithm is used because of the following advantages over CPU-bound and Memory-bound algorithms:

- **Computation guarantee**. The computation guarantee means a cryptographic puzzle guarantees a lower and upper bound on the number of cryptographic operations spent on a client to find the puzzle answer. In other words, a mali- cious client should not be able to solve a puzzle spending significantly less number of operations than required.
- **Efficiency**. The construction, distribution and verification of a puzzle by the server should be efficient, in terms of CPU, memory, bandwidth, hard disk etc. Specifically, puzzle construction, distribution and verification should add minimal overhead to the server to prevent the puzzle scheme from becoming an avenue for denying service.
- **Adjustability of difficulty**. This property is also referred to as puzzle granularity. Adjustability of puzzle difficulty means the cost of solving the puzzle can be increased or decreased in fine granularity. Adjustability of difficulty is important, because finer adjustability enables the server to achieve better trade-off between blocking attackers and the service degradation of legitimate clients.
- **Correlation-free**. A puzzle is considered correlation-free if knowing the solutions to all previous puzzles seen by a client does not make solving a new puzzle any easy. Apparently, if a puzzle is not correlation-free, then it allows malicious clients to solve puzzles faster by correlating previous answers.
- **Stateless**. A puzzle is said to be stateless if it requires the server to store no client information or puzzle-related data in order to verify puzzle solutions. Requiring the server to use a small and fixed memory for storing such information is also acceptable in most cases.
- **Tamper-resistance**. A puzzle scheme should limit replay attacks over time and space. Puzzle solutions should not be valid indefinitely and should not be usable by other clients.
- **Non-parallelizability**. Non-parallelizability means the puzzle solution cannot be computed in parallel using multiple machines. Non-parallelizable puzzles can prevent attackers from distributing computation of a puzzle solution to a group of machines to obtain the solution quicker.
- **Puzzle fairness**. Puzzle fairness means a puzzle should take same amount of time for all clients to compute, re- gardless of their CPU power, memory size, and bandwidth. If a puzzle can achieve fairness, then a powerful DoS attacker can effectively be reduced to a legitimate client. Not being able to achieve fairness leads to the resource disparity problem we mentioned earlier.
- **Minimum interference**. This property requires that puz- zle computation at the client should not interfere with user’s normal operations. If a puzzle scheme takes up too much resources and interfere with users’ normal computing activity, users might disable the puzzle scheme or even try to avoid using any service that deploys such a puzzle scheme.

Reference: [A Guided Tour Puzzle for Denial of Service Prevention](https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.597.6304&rep=rep1&type=pdf)

## Run demo

Following command starts server, 2 guide and client instances and outputs debug information about client-server-guide communication:
```shell
docker-compose up
```