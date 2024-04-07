"use server";

async function createThreadChatAPI(accessToken, body = {}) {
  try {
    const response = await fetch(`${process.env.ZYG_API_URL}/-/threads/chat/`, {
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
      `${process.env.ZYG_API_URL}/-/threads/chat/${threadId}/messages/`,
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

export async function createThreadChat(values) {
  try {
    // customerId: c_co61abktidu1t3i3dn60
    const jwt = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ3b3Jrc3BhY2VJZCI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwIiwiZXh0ZXJuYWxJZCI6Inh4eHgtMTExLXp6enoiLCJlbWFpbCI6InNhbmNoaXRycmtAZ21haWwuY29tIiwicGhvbmUiOiIrOTE3NzYwNjg2MDY4IiwiaXNzIjoiYXV0aC56eWcuYWkiLCJzdWIiOiJjX2NvNjFhYmt0aWR1MXQzaTNkbjYwIiwiYXVkIjpbImN1c3RvbWVyIl0sImV4cCI6MTc0Mzc1Nzg3MSwibmJmIjoxNzEyMjIxODcxLCJpYXQiOjE3MTIyMjE4NzEsImp0aSI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwOmNfY282MWFia3RpZHUxdDNpM2RuNjAifQ.epCQ4aXvYPXIhVrX6TtfYrq0XxYXT18kIWsOae8HvUQ`;
    const { ...rest } = values;
    const [err, data] = await createThreadChatAPI(jwt, rest);
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

export async function sendThreadChatMessage(threadId, values) {
  try {
    // customerId: c_co61abktidu1t3i3dn60
    const jwt = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ3b3Jrc3BhY2VJZCI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwIiwiZXh0ZXJuYWxJZCI6Inh4eHgtMTExLXp6enoiLCJlbWFpbCI6InNhbmNoaXRycmtAZ21haWwuY29tIiwicGhvbmUiOiIrOTE3NzYwNjg2MDY4IiwiaXNzIjoiYXV0aC56eWcuYWkiLCJzdWIiOiJjX2NvNjFhYmt0aWR1MXQzaTNkbjYwIiwiYXVkIjpbImN1c3RvbWVyIl0sImV4cCI6MTc0Mzc1Nzg3MSwibmJmIjoxNzEyMjIxODcxLCJpYXQiOjE3MTIyMjE4NzEsImp0aSI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwOmNfY282MWFia3RpZHUxdDNpM2RuNjAifQ.epCQ4aXvYPXIhVrX6TtfYrq0XxYXT18kIWsOae8HvUQ`;
    const { ...rest } = values;
    const [err, data] = await sendThreadChatMessageAPI(jwt, threadId, rest);
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
