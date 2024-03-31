"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { createClient } from "@/utils/supabase/actions";
import { getAuthToken } from "@/utils/supabase/helpers";

async function createWorkspaceAPI(accessToken, body = {}) {
  try {
    const response = await fetch(`${process.env.ZYG_API_URL}/workspaces/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify({ ...body }),
    });
    if (!response.ok) {
      const { status, statusText } = response;
      return [
        new Error(
          `error creating workspace with status: ${status} and statusText: ${statusText}`
        ),
        null,
      ];
    }
    const data = await response.json();
    const { slug } = data;
    console.log(`successfully created workspace with slug: ${slug}`);
    return [null, { slug }];
  } catch (err) {
    return [err, null];
  }
}

/**
 * Creates a workspace.
 * @param {Object} values - The values for creating the workspace.
 * @returns {Promise<Object>} - A promise that resolves to an object containing the error and data.
 */
export async function createWorkspace(values) {
  const cookieStore = cookies();
  const supabase = createClient(cookieStore);
  try {
    const accessToken = await getAuthToken(supabase);
    const [err, workspace] = await createWorkspaceAPI(accessToken, values);
    if (err) {
      return {
        error: {
          message: "Workspace creation failed",
        },
        data: null,
      };
    }
    revalidatePath("/workspaces/");
    return {
      error: null,
      data: workspace,
    };
  } catch (err) {
    console.error(err);
    return {
      error: {
        message: "Something went wrong!",
      },
      data: null,
    };
  }
}
