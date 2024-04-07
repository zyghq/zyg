"use client";
import * as React from "react";
import usePartySocket from "partysocket/react";
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { SendHorizonalIcon } from "lucide-react";
import { nanoid } from "nanoid";

const PARTY_KIT_HOST = "127.0.0.1:1999";

const identity = async (socket) => {
  console.log("do some auth here....");
  // checkout the nextjs chat sample for example, on
  // how is the _pk is being used.
  console.log("socket._pk: ", socket._pk);
};

function Message({ message }) {
  return (
    <div className="flex space-x-2">
      <Avatar className="h-6 w-6">
        <AvatarImage alt="User" src="/images/profile.jpg" />
        <AvatarFallback>U</AvatarFallback>
      </Avatar>
      <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
        <p className="text-sm">{message.text}</p>
      </div>
    </div>
  );
}

// we can load initial messages via GET request
export default function Room({ roomId, messages: initialMessages = [] }) {
  const [messages, setMessages] = React.useState(initialMessages);
  const socket = usePartySocket({
    host: PARTY_KIT_HOST,
    room: roomId,
    onOpen: (e) => {
      if (e.target) {
        console.log("connected to room", roomId);
        identity(e.target);
      }
    },
    onMessage: (event) => {
      console.log("message event", event);
      const { data } = event;
      const message = JSON.parse(data);
      setMessages((prev) => [...prev, message]);
    },
    onClose: () => console.log("disconnected from room", roomId),
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    const text = e.currentTarget.message.value;
    if (text?.trim()) {
      const message = {
        id: nanoid(),
        text,
        from: {
          id: "rahul",
        },
        at: Math.floor(new Date().getTime() / 1000),
      };
      socket.send(JSON.stringify(message));
      e.currentTarget.message.value = "";
      scrollToBottom();
    }
  };

  const onEnterPress = (e) => {
    if (e.keyCode === 13 && e.shiftKey === false) {
      e.preventDefault();
      e.target.form.requestSubmit();
    }
  };

  const scrollToBottom = () => {};

  React.useEffect(() => {
    console.log("entered room with id", roomId);
  }, [roomId]);

  return (
    <React.Fragment>
      <ScrollArea className="p-4 h-[calc(100dvh-2rem)]">
        <div className="space-y-2">
          {messages.map((message) => (
            <Message key={message.id} message={message} />
          ))}
          {/* <div className="max-w-xs mr-auto">
            <div className="flex space-x-2">
              <Avatar className="h-6 w-6">
                <AvatarImage alt="User" src="/images/profile.jpg" />
                <AvatarFallback>U</AvatarFallback>
              </Avatar>
              <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                <p className="text-sm">
                  {
                    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ac nunc auctor, lacinia nisl id, efficitur nunc. Nulla facilisi. Sed euismod, nisl nec ultrices ultricies, nunc nunc aliquet nunc, sit amet tincidunt nunc nunc eu nunc. Nulla facilisi. Sed euismod, nisl nec ultrices ultricies, nunc nunc aliquet nunc, sit amet tincidunt nunc nunc eu nunc."
                  }
                </p>
              </div>
            </div>
          </div> */}
          {/* <div className="max-w-xs mr-auto">
            <div className="flex space-x-2">
              <Avatar className="h-6 w-6">
                <AvatarImage alt="User" src="/images/profile.jpg" />
                <AvatarFallback>U</AvatarFallback>
              </Avatar>
              <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                <p className="text-sm">
                  {
                    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ac nunc auctor, lacinia nisl id, efficitur nunc. Nulla facilisi. Sed euismod, nisl nec ultrices ultricies, nunc nunc aliquet nunc, sit amet tincidunt nunc nunc eu nunc. Nulla facilisi. Sed euismod, nisl nec ultrices ultricies, nunc nunc aliquet nunc, sit amet tincidunt nunc nunc eu nunc."
                  }
                </p>
              </div>
            </div>
          </div> */}
          {/* <div className="max-w-xs ml-auto">
            <div className="flex space-x-2">
              <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                <p className="text-sm">
                  {
                    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ac nunc auctor, lacinia nisl id, efficitur nunc. Nulla facilisi. Sed euismod, nisl nec ultrices ultricies, nunc nunc aliquet nunc, sit amet tincidunt nunc nunc eu nunc. Nulla facilisi. Sed euismod, nisl nec ultrices ultricies, nunc nunc aliquet nunc, sit amet tincidunt nunc nunc eu nunc."
                  }
                </p>
              </div>
              <Avatar className="h-6 w-6">
                <AvatarImage alt="User" src="/images/profile.jpg" />
                <AvatarFallback>U</AvatarFallback>
              </Avatar>
            </div>
          </div> */}
          {/* <div ref={messageEndRef} /> */}
        </div>
      </ScrollArea>
      <form onSubmit={handleSubmit} className="flex items-center p-2 mt-auto">
        <Textarea
          type="text"
          name="message"
          placeholder="Send us a message"
          title="Send us a message"
          className="mr-1"
          rows={1}
          onKeyDown={onEnterPress}
        />
        <Button size="icon" type="submit">
          <SendHorizonalIcon className="h-4 w-4" />
        </Button>
      </form>
    </React.Fragment>
  );
}
