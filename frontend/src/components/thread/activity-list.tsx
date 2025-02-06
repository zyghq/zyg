import type React from "react"

import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { useState } from "react"

export function EmailThread() {
  const [selectedCard, setSelectedCard] = useState<null | number>(null)

  const emails = [
    {
      sender: "Sanchit Rk",
      recipient: "....@inbound.postmark...",
      date: "9/30/2024",
      content: (
        <>
          <p>Hi there,</p>
          <p>
            I think the new deployment is having an issue. There are logs in AWS that seem to have errors. Can we check
            on that? I can see them in the AWS logs.
          </p>
          <div className="space-y-1">
            <p>Thanks,</p>
            <p>Sanchit Rk</p>
          </div>
        </>
      ),
      hasActions: true,
    },
    {
      sender: "Sanchit at Zyg",
      recipient: "Sanchit Rk",
      date: "10/28/2024",
      content: (
        <>
          <p>Yeah we been working on it, it seems it there are some deployment issues.</p>
          <div className="pl-2 sm:pl-4 border-l-2 border-muted text-muted-foreground space-y-2 text-sm sm:text-base">
            <p className="text-xs sm:text-sm">
              On Mon, 30 Sep 2024 at 07:40, Sanchit Rk{" "}
              <a className="text-primary dark:text-indigo-400" href="mailto:sanchitrk@gmail.com">
                sanchitrk@gmail.com
              </a>{" "}
              wrote:
            </p>
            <p>Hi there,</p>
            <p>
              I think the new deployment is having an issue. There are logs in AWS that seem to have errors. Can we
              check on that? I can see them in the AWS logs.
            </p>
            <div>
              <p>Thanks,</p>
              <p>Sanchit Rk</p>
            </div>
          </div>
        </>
      ),
      hasActions: true,
    },
    {
      sender: "Sanchit at Zyg",
      recipient: "Sanchit Rk",
      date: "10/28/2024",
      content: <p>This issue should be fixed soon</p>,
      hasActions: false,
    },
    {
      sender: "Sanchit Rk",
      recipient: "Sanchit at Zyg",
      date: "10/28/2024",
      content: <p>Thanks for the update, Zyg. I'm confident the team will get this resolved soon.</p>,
      hasActions: false,
    },
    {
      sender: "Sanchit Rk",
      recipient: "....@inbound.postmarkapp.com",
      date: "11/5/2024",
      content: <p>We do not have any urgency for this issue</p>,
      hasActions: false,
    },
  ]

  return (
    <div className="max-w-3xl mx-auto space-y-4 p-2 sm:p-4">
      {emails.map((email, index) => (
        <Card
          className={`shadow-sm hover:shadow-md transition-shadow relative group cursor-pointer ${
            selectedCard === index ? "shadow-md dark:shadow-indigo-500/20" : ""
          } dark:bg-gray-800 dark:hover:shadow-indigo-500/10`}
          key={index}
          onClick={() => setSelectedCard(selectedCard === index ? null : index)}
        >
          <div
            className={`absolute inset-0 border rounded-lg ${
              selectedCard === index ? "opacity-100" : "opacity-0"
            } pointer-events-none border-indigo-200 dark:border-indigo-500`}
          />
          <CardHeader className="flex flex-col sm:flex-row items-start sm:items-center justify-between space-y-2 sm:space-y-0 pb-2 dark:text-gray-200">
            <div className="flex flex-col sm:flex-row sm:items-center gap-1">
              <span className="text-xs sm:text-sm font-medium truncate max-w-[200px] sm:max-w-none">
                {email.sender} to {email.recipient}
              </span>
              <span className="text-xs text-muted-foreground dark:text-gray-400 sm:ml-2">Shared</span>
            </div>
            <div className="flex items-center gap-2">
              {email.hasActions && (
                <>
                  <button className="p-1 hover:bg-accent rounded-md">
                    <ReplyIcon className="h-4 w-4 text-muted-foreground dark:text-gray-400" />
                  </button>
                  <button className="p-1 hover:bg-accent rounded-md">
                    <ForwardIcon className="h-4 w-4 text-muted-foreground dark:text-gray-400" />
                  </button>
                </>
              )}
              <span className="text-xs text-muted-foreground dark:text-gray-400">{email.date}</span>
            </div>
          </CardHeader>
          <CardContent className="space-y-2 dark:text-gray-300 text-sm sm:text-base">{email.content}</CardContent>
        </Card>
      ))}
    </div>
  )
}

function ReplyIcon(props: React.ComponentProps<"svg">) {
  return (
    <svg
      fill="none"
      height="24"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      viewBox="0 0 24 24"
      width="24"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <polyline points="9 17 4 12 9 7" />
      <path d="M20 18v-2a4 4 0 0 0-4-4H4" />
    </svg>
  )
}

function ForwardIcon(props: React.ComponentProps<"svg">) {
  return (
    <svg
      fill="none"
      height="24"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      viewBox="0 0 24 24"
      width="24"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <polyline points="15 17 20 12 15 7" />
      <path d="M4 18v-2a4 4 0 0 1 4-4h12" />
    </svg>
  )
}
