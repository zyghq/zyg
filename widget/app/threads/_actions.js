"use server";

async function createThreadChatAPI(accessToken, body = {}) {
  try {
    const response = await fetch(`${process.env.ZYG_XAPI_URL}/threads/chat/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify(body),
    });
    if (!response.ok) {
      const { status, statusText } = response;
      return [
        new Error(
          `error creating thread chat with status: ${status} - ${statusText}`
        ),
        null,
      ];
    }
    const data = await response.json();
    return [null, data];
  } catch (err) {
    return [err, null];
  }
}

async function sendThreadChatMessageAPI(accessToken, threadId, body = {}) {
  try {
    const { message } = body;
    const response = await fetch(
      `${process.env.ZYG_XAPI_URL}/threads/chat/${threadId}/messages/`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          message,
        }),
      }
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return [
        new Error(
          `error creating thread chat with status: ${status} - ${statusText}`
        ),
        null,
      ];
    }
    const data = await response.json();
    return [null, data];
  } catch (err) {
    return [err, null];
  }
}

export async function createThreadChat(authToken, values) {
  try {
    const [err, data] = await createThreadChatAPI(authToken, { ...values });
    if (err) {
      console.error(err);
      return {
        error: {
          message: "Failed. Please try again later.",
        },
        data: null,
      };
    }
    // TODO: do revalidate cache for list of threads.
    return {
      error: null,
      data,
    };
  } catch (err) {
    console.error(err);
    return {
      error: {
        message: "Failed. Please try againg later.",
      },
      data: null,
    };
  }
}

export async function sendThreadChatMessage(authToken, threadId, values) {
  console.log(
    `invoking sendThreadChatMessage with authToken: ${authToken} and threadId: ${threadId}`
  );
  try {
    const [err, data] = await sendThreadChatMessageAPI(authToken, threadId, {
      ...values,
    });
    if (err) {
      console.error(err);
      return {
        error: {
          message: "Failed. Please try again later.",
        },
        data: null,
      };
    }
    // TODO: do revalidate cache for list of threads.
    return {
      error: null,
      data,
    };
  } catch (err) {
    console.error(err);
    return {
      error: {
        message: "Failed. Please try againg later.",
      },
      data: null,
    };
  }
}
