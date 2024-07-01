import { NextResponse } from "next/server";

const ZYG_AUTH_COOKIE_NAME = "__zygtoken";

export async function GET(request, { params }) {
  const { threadId } = params;
  const token = request.cookies.get(ZYG_AUTH_COOKIE_NAME);
  if (!token) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }
  const { value } = token;
  try {
    const resp = await fetch(
      `${process.env.NEXT_PUBLIC_XAPI_URL}/threads/chat/${threadId}/messages/`,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${value}`,
        },
      }
    );

    if (!resp.ok) {
      return NextResponse.json(
        { error: "authentication error" },
        { status: 401 }
      );
    }
    const data = await resp.json();
    return NextResponse.json({ ...data }, { status: 200 });
  } catch (err) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }
}
