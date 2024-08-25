"use server";

import { revalidatePath } from "next/cache";
import { z } from "zod";
import { createThreadResponseSchema, CreateThreadResponse } from "@/lib/thread";

interface CreateThreadBody {
  message: string;
}

interface UpdateEmailBody {
  email: string;
}

export async function createThreadActionAPI(
  widgetId: string,
  jwt: string,
  body: CreateThreadBody
): Promise<{
  error: { message: string } | null;
  data: CreateThreadResponse | null;
}> {
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
    try {
      const thread = createThreadResponseSchema.parse(data);
      return {
        error: null,
        data: thread,
      };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(
          "Failed response schema validation, update threadResponseSchema"
        );
        console.error(err.message);
      } else console.error(err);
      return {
        error: {
          message: "Failed response schema validation",
        },
        data: null,
      };
    }
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

// required, as we cannot have try..catch due to next.js revalidatePath.
export async function createThreadAction(
  widgetId: string,
  jwt: string,
  body: CreateThreadBody
): Promise<{
  error: { message: string } | null;
  data: CreateThreadResponse | null;
}> {
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

export async function updateEmailActionAPI(
  widgetId: string,
  jwt: string,
  body: UpdateEmailBody
) {
  try {
    const response = await fetch(
      `${process.env.ZYG_XAPI_URL}/widgets/${widgetId}/me/identities/`,
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
      console.error(`Failed to update email. Status: ${status} ${statusText}`);
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

export async function updateEmailAction(
  widgetId: string,
  jwt: string,
  body: UpdateEmailBody
) {
  const { error, data } = await updateEmailActionAPI(widgetId, jwt, body);
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
