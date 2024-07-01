import { NextResponse } from "next/server";
import { cookies } from "next/headers";

const ZYG_AUTH_COOKIE_NAME = "__zygtoken";

export async function POST(request) {
  const { token } = await request.json();
  console.log("URL ->", process.env.NEXT_PUBLIC_XAPI_URL);
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

    console.log("In POST /me/ Data ->", data);

    const response = NextResponse.json({ ...data, authToken }, { status: 200 });
    response.cookies.set(ZYG_AUTH_COOKIE_NAME, token, {
      httpOnly: false, // make sure it is accessible by the browser (client)
      // secure: process.env.NODE_ENV === "production",
      secure: true,
      sameSite: "None",
      maxAge: 60 * 60 * 24 * 7,
      path: "/",
    });
    return response;
  } catch (err) {
    console.log("In POST /me/ Error ->", err);
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }
}

export async function GET(request) {
  const token = request.cookies.get(ZYG_AUTH_COOKIE_NAME);
  if (!token) {
    console.log("In GET /me/ No Token ->", token);
    return NextResponse.json(
      { error: "authentication error" },
      { status: 401 }
    );
  }

  console.log("In GET /me/ got Token...");
  const { value } = token;

  console.log("In GET /me/ Token ->", token);
  try {
    const resp = await fetch(`${process.env.NEXT_PUBLIC_XAPI_URL}/me/`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${value}`,
      },
    });

    const { status, statusText } = resp;
    console.log("In GET /me/ Status ->", status);
    console.log("In GET /me/ Status Text ->", statusText);

    if (!resp.ok) {
      console.log("In GET NOT OK /me/ Error ->", resp);
      return NextResponse.json(
        { error: "authentication error" },
        { status: 401 }
      );
    }
    const data = await resp.json();
    console.log("In GET /me/ Data ->", data);
    return NextResponse.json({ ...data, authToken: token }, { status: 200 });
  } catch (err) {
    console.log("In GET /me/ Error ->", err);
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
