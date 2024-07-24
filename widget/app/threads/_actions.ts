"use server";

import { revalidatePath } from "next/cache";

interface CreateThreadBody {
  message: string;
}

export async function createThreadActionAPI(
  widgetId: string,
  jwt: string,
  body: CreateThreadBody
) {
  try {
    const response = await fetch(
      `${process.env.ZYG_XAPI_URL}/widgets/${widgetId}/threads/chat/`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
        body: JSON.stringify(body),
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      console.error(`Failed to create thread. Status: ${status} ${statusText}`);
      return {
        data: null,
        error: {
          message: "Failed. Please try again later.",
        },
      };
    }
    const data = await response.json();
    return {
      error: null,
      data,
    };
  } catch (err) {
    console.error("Something went wrong", err);
    return {
      error: {
        message: "Something went wrong. Please try again later.",
      },
      data: null,
    };
  }
}

export async function createThreadAction(
  widgetId: string,
  jwt: string,
  body: CreateThreadBody
) {
  const { error, data } = await createThreadActionAPI(widgetId, jwt, body);
  if (error) {
    return {
      error,
      data: null,
    };
  }
  revalidatePath("/");
  return {
    error: null,
    data,
  };
}

interface SendMessageBody {
  message: string;
}

export async function sendThreadMessageActionAPI(
  widgetId: string,
  threadId: string,
  jwt: string,
  body: SendMessageBody
) {
  try {
    const response = await fetch(
      `${process.env.ZYG_XAPI_URL}/widgets/${widgetId}/threads/chat/${threadId}/messages/`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
        body: JSON.stringify(body),
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      console.error(
        `Failed to send thread message. Status: ${status} ${statusText}`
      );
      return {
        data: null,
        error: {
          message: "Failed. Please try again later.",
        },
      };
    }
    const data = await response.json();
    return {
      error: null,
      data,
    };
  } catch (err) {
    console.error("Something went wrong", err);
    return {
      error: {
        message: "Something went wrong. Please try again later.",
      },
      data: null,
    };
  }
}

export async function sendThreadMessageAction(
  widgetId: string,
  threadId: string,
  jwt: string,
  body: SendMessageBody
) {
  const { error, data } = await sendThreadMessageActionAPI(
    widgetId,
    threadId,
    jwt,
    body
  );
  if (error) {
    return {
      error,
      data: null,
    };
  }
  revalidatePath("/");
  return {
    error: null,
    data,
  };
}
