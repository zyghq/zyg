import { NextResponse } from "next/server";
import { cookies } from "next/headers";

const ZYG_AUTH_COOKIE_NAME = "__zygtoken";

export async function POST(request) {
  const { token } = await request.json();
  try {
    const resp = await fetch(`${process.env.NEXT_PUBLIC_XAPI_URL}/me/`, {
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
    const authToken = {
      value: token,
      name: ZYG_AUTH_COOKIE_NAME,
    };

    const response = NextResponse.json({ ...data, authToken }, { status: 200 });
    response.cookies.set(ZYG_AUTH_COOKIE_NAME, token, {
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
  const token = request.cookies.get(ZYG_AUTH_COOKIE_NAME);
  if (!token) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }

  const { value } = token;
  try {
    const resp = await fetch(`${process.env.NEXT_PUBLIC_XAPI_URL}/me/`, {
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
    return NextResponse.json({ ...data, authToken: token }, { status: 200 });
  } catch (err) {
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }
}

export async function DELETE() {
  const cookieStore = cookies();
  cookieStore.delete(ZYG_AUTH_COOKIE_NAME);
  return NextResponse.json({ ok: true }, { status: 200 });
}
