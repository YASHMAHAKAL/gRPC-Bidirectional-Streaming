import grpc
import threading
import time

# Import the classes generated from your .proto file
import chat_pb2
import chat_pb2_grpc

def listen_for_messages(stub, response_iterator):
    """
    This function runs in a separate thread, continuously listening for messages
    from the server and printing them to the console.
    """
    try:
        # The 'for' loop will block and wait for new messages from the server
        for message in response_iterator:
            # When a message is received, print it and then re-print the input prompt
            print(f"\n[{message.user}]: {message.message.strip()}")
            print("You: ", end="", flush=True)
    except grpc.RpcError as e:
        # Handle cases where the server connection is lost
        print(f"\nConnection lost. Error receiving message: {e}")

def run():
    """
    Main function to set up the gRPC channel and manage the chat session.
    """
    # Establish a connection to the gRPC server (your existing Go server)
    with grpc.insecure_channel('localhost:50053') as channel:
        # Create a client "stub"
        stub = chat_pb2_grpc.ChatServiceStub(channel)

        user = input("Enter your name: ")

        def message_iterator():
            """
            An iterator that yields messages from the user's input.
            This is passed to the ChatStream method.
            """
            while True:
                text = input("You: ")
                message = chat_pb2.ChatMessage(
                    user=user,
                    message=text,
                    timestamp=int(time.time())
                )
                yield message

        # Call the ChatStream RPC. This returns an iterator for server responses.
        response_iterator = stub.ChatStream(message_iterator())

        # Start a new thread to listen for incoming messages from the server.
        # This allows you to type messages while simultaneously receiving them.
        threading.Thread(target=listen_for_messages, args=(stub, response_iterator), daemon=True).start()

        # Keep the main thread alive to allow the background thread to continue running.
        # The program will exit if the main thread finishes.
        try:
            while threading.main_thread().is_alive():
                time.sleep(1)
        except KeyboardInterrupt:
            print("\nExiting chat.")

if __name__ == '__main__':
    run()
