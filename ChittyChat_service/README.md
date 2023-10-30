# ChittyChat

ChittyChat is a real-time chat service built using gRPC, allowing users to communicate in a chatroom-like environment via CLI. It provides bidirectional message streaming between clients and the server, enabling a continuous exchange of messages.


### How to launch this ChittyChat:

1. "go run ChittyChatServer.go" to launch the server application. It will listen on port 5001. <br>
2. For each Client you wish to use, open a new terminal and write "go run ChittyChatClient.go" to launch the client application. <br>
3. Provide the Client application with a UTF8 valid username(No spaces), followed by pressing enter to register the client. <br>
4. The client will now connect to the server application, and you are now a part of the chat.<br>
5. To send a message, simply type your message in your terminal and press enter to send.
6. To terminate a client press "control + c" or simply close the relevant terminal.

