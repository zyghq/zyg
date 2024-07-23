import { NextRequest, NextResponse } from "next/server";

export async function POST(
  request: NextRequest,
  { params }: { params: { widgetId: string; threadId: string } }
) {
  try {
    const { widgetId, threadId } = params;
    const body = await request.json();
    const { jwt } = body;
    const response = await fetch(
      `${process.env.ZYG_XAPI_URL}/widgets/${widgetId}/threads/chat/${threadId}/messages/`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
      }
    );

    if (!response.ok) {
      return NextResponse.json({ error: "Bad Request" }, { status: 400 });
    }

    const data = await response.json();
    return NextResponse.json(data, { status: 200 });
  } catch (err) {
    console.error("Error processing request:", err);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
