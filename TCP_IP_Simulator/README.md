                                        -Group: DSGO - Handin2-

Members: Bekr, Mbek og Olau. 



*A) What are packages in your implementation? What data structure do you use to transmit data and meta-data?*
- We use channels to transmit data in the form of integers.
MAY BE UPDATED TO TRANSMIT STRUCTS HOLDING DIFFERENT BOOLS AND INTS AND STRINGS


*B) Does your implementation use threads or processes? Why is it not realistic to use threads?*
- This implementation uses threads. 
WHY IS IT NOT REALISTIC TO USE THREADS?


*C) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?*
- We would use the sequence number to put the messages back in order. 
We have read that TCP requests a re-transmission if enough packets are lost,
however this seems painfully overkill to implement in our code.


*D) In case messages can be delayed or lost, how does your implementation handle message loss?*
- It does not...
Message delay is only an issue if it messes with the message order, or if they are infinitely long.
Otherwise hosts will wait indefinitely for messages.
We might have to request a retransmission to recover lost messages.


*E) Why is the 3-way handshake important?*
- The TCP/IP protocol is a connection-oriented protocol. Without a 3-way handshake, 
we cannot confirm that we have a stable and secure connection to stream data on, 
which is one of the main points of a TCP/IP connection.