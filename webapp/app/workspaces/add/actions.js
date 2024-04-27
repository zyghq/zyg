"use server";

import { revalidatePath } from "next/cache";
import { createClient } from "@/utils/supabase/server";
import { getAuthToken } from "@/utils/supabase/helpers";

async function createWorkspaceAPI(accessToken, body = {}) {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({ ...body }),
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return [
        new Error(
          `error creating workspace with status: ${status} and statusText: ${statusText}`,
        ),
        null,
      ];
    }
    const data = await response.json();
    const { workspaceId } = data;
    console.log(
      `successfully created workspace with workspaceId: ${workspaceId}`,
    );
    return [null, { workspaceId }];
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
  const supabase = createClient();
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
