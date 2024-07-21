import { NextRequest, NextResponse } from "next/server";

export async function POST(
  request: NextRequest,
  { params }: { params: { widgetId: string } }
) {
  try {
    const { widgetId } = params;
    const body = await request.json();
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_XAPI_URL}/widgets/${widgetId}/init/`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
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
