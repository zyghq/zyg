import { Button } from "@/components/ui/button";
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "./ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { DotFilledIcon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Header } from "@/components/headers";

//  lets move some of the components to the Page.jsx component.
export default function Widget() {
  return (
    <div className="flex flex-col h-full w-full max-w-md rounded-lg overflow-hidden">
      <Header label="Hey! How can we help?" />
      <div className="px-4">
        <Button variant="secondary" className="w-full" asChild>
          <Link href="/search/" className="flex">
            <MagnifyingGlassIcon className="h-4 w-4 mr-1" />
            Search
          </Link>
        </Button>
      </div>
      <div className="flex-1 p-4">
        <Tabs defaultValue="home">
          <TabsList>
            <TabsTrigger value="home">Home</TabsTrigger>
            <TabsTrigger value="threads">Threads</TabsTrigger>
          </TabsList>
          <TabsContent value="home">
            <ScrollArea className="h-96">
              <div className="space-y-2">
                <div className="flex space-x-2">
                  <Avatar className="h-6 w-6">
                    <AvatarImage alt="User" src="/images/profile.jpg" />
                    <AvatarFallback>U</AvatarFallback>
                  </Avatar>
                  <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                    <p className="text-sm">Hello, how can I help you?</p>
                  </div>
                </div>
                <div className="flex space-x-2">
                  <Avatar className="h-6 w-6">
                    <AvatarImage alt="User" src="/images/profile.jpg" />
                    <AvatarFallback>U</AvatarFallback>
                  </Avatar>
                  <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                    <div className="text-sm">
                      {`Let us use the space in Home for more quick feeds, anouncements, and updates. Recent active Threads, etc.`}
                    </div>
                  </div>
                </div>
              </div>
            </ScrollArea>
          </TabsContent>
          <TabsContent value="threads">
            <ScrollArea className="flex-1 pb-2">
              <Card>
                <CardHeader>
                  <CardTitle>
                    Not able to login to the app. Need help!
                  </CardTitle>
                  <CardDescription>
                    <div className="flex justify-between pt-2">
                      <div className="flex">
                        <Avatar className="h-6 w-6">
                          <AvatarImage alt="User" src="/images/profile.jpg" />
                          <AvatarFallback>T</AvatarFallback>
                        </Avatar>
                        <div className="px-1 font-semibold">Tom Barr</div>
                      </div>
                      <div className="px-1 text-sm">1 hr. ago</div>
                    </div>
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-sm">
                    This issue is now resolved. Closing now.
                  </p>
                </CardContent>
                <CardFooter className="justify-between">
                  <Badge variant="outlined">
                    <DotFilledIcon className="text-green-500" />
                    Done
                  </Badge>
                </CardFooter>
              </Card>
            </ScrollArea>
          </TabsContent>
        </Tabs>
      </div>
      <div className="p-4 border-t">
        <form className="flex space-x-2">
          <Input
            className="flex-1"
            placeholder="Type your message here"
            type="text"
          />
          <Button type="submit">
            <SendIcon className="h-6 w-6" />
          </Button>
        </form>
      </div>
    </div>
  );
}

// function FlagIcon(props) {
//   return (
//     <svg
//       {...props}
//       xmlns="http://www.w3.org/2000/svg"
//       width="24"
//       height="24"
//       viewBox="0 0 24 24"
//       fill="none"
//       stroke="currentColor"
//       strokeWidth="2"
//       strokeLinecap="round"
//       strokeLinejoin="round"
//     >
//       <path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z" />
//       <line x1="4" x2="4" y1="22" y2="15" />
//     </svg>
//   );
// }

// function XIcon(props) {
//   return (
//     <svg
//       {...props}
//       xmlns="http://www.w3.org/2000/svg"
//       width="24"
//       height="24"
//       viewBox="0 0 24 24"
//       fill="none"
//       stroke="currentColor"
//       strokeWidth="2"
//       strokeLinecap="round"
//       strokeLinejoin="round"
//     >
//       <path d="M18 6 6 18" />
//       <path d="m6 6 12 12" />
//     </svg>
//   );
// }

function SendIcon(props) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="m22 2-7 20-4-9-9-4Z" />
      <path d="M22 2 11 13" />
    </svg>
  );
}

// function HomeIcon(props) {
//   return (
//     <svg
//       {...props}
//       xmlns="http://www.w3.org/2000/svg"
//       width="24"
//       height="24"
//       viewBox="0 0 24 24"
//       fill="none"
//       stroke="currentColor"
//       strokeWidth="2"
//       strokeLinecap="round"
//       strokeLinejoin="round"
//     >
//       <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
//       <polyline points="9 22 9 12 15 12 15 22" />
//     </svg>
//   );
// }

// function UserIcon(props) {
//   return (
//     <svg
//       {...props}
//       xmlns="http://www.w3.org/2000/svg"
//       width="24"
//       height="24"
//       viewBox="0 0 24 24"
//       fill="none"
//       stroke="currentColor"
//       strokeWidth="2"
//       strokeLinecap="round"
//       strokeLinejoin="round"
//     >
//       <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
//       <circle cx="12" cy="7" r="4" />
//     </svg>
//   );
// }

// function SettingsIcon(props) {
//   return (
//     <svg
//       {...props}
//       xmlns="http://www.w3.org/2000/svg"
//       width="24"
//       height="24"
//       viewBox="0 0 24 24"
//       fill="none"
//       stroke="currentColor"
//       strokeWidth="2"
//       strokeLinecap="round"
//       strokeLinejoin="round"
//     >
//       <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
//       <circle cx="12" cy="12" r="3" />
//     </svg>
//   );
// }
