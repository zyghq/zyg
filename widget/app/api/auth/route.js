import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function POST(request) {
  const { token } = await request.json();
  try {
    const resp = await fetch(`${process.env.ZYG_API_URL}/-/me/`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (!resp.ok) {
      return NextResponse.json(
        { error: "authentication error" },
        { status: 401 }
      );
    }
    const data = await resp.json();

    const response = NextResponse.json({ ...data }, { status: 200 });
    response.cookies.set("__zygtoken", token, {
      httpOnly: false, // make sure it is accessible by the browser (client)
      secure: process.env.NODE_ENV === "production",
      sameSite: "strict",
      maxAge: 60 * 60 * 24 * 7,
      path: "/",
    });
    return response;
  } catch (err) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }
}

export async function GET(request) {
  const token = request.cookies.get("__zygtoken");
  if (!token) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }

  const { value } = token;
  try {
    const resp = await fetch(`${process.env.ZYG_API_URL}/-/me/`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${value}`,
      },
    });

    if (!resp.ok) {
      return NextResponse.json(
        { error: "authentication error" },
        { status: 401 }
      );
    }
    const data = await resp.json();
    console.log("response from the ZYG API....", data);

    return NextResponse.json({ ...data }, { status: 200 });
  } catch (err) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }
}

export async function DELETE() {
  const cookieStore = cookies();
  cookieStore.delete("__zygtoken");
  return NextResponse.json({ ok: true }, { status: 200 });
}
