import { z } from "zod";
import { fallback } from "@tanstack/router-zod-adapter";
//
// for more: https://tanstack.com/router/latest/docs/framework/react/guide/search-params
// usage of `.catch` or `default` matters.
const stagesSchema = (validValues: string[]) => {
  const sanitizeArray = (arr: string[]) => {
    // remove duplicates
    const uniqueValues = [...new Set(arr)];
    // filter only valid values
    const uniqueValidValues: string[] = uniqueValues.filter((val) =>
      validValues.includes(val)
    );

    // no valid values
    if (uniqueValidValues.length === 0) {
      throw new Error("invalid statuses passed");
    }

    if (uniqueValidValues.length === 1) {
      return uniqueValidValues[0];
    }

    return uniqueValidValues;
  };
  return z.union([
    z.string().refine((value) => validValues.includes(value)),
    z.array(z.string()).transform(sanitizeArray),
    z.undefined(),
  ]);
};

const prioritiesSchema = (validValues: string[]) => {
  const sanitizeArray = (arr: string[]) => {
    // remove duplicates
    const uniqueValues = [...new Set(arr)];
    // filter only valid values
    const uniqueValidValues: string[] = uniqueValues.filter((val) =>
      validValues.includes(val)
    );

    // no valid values
    if (uniqueValidValues.length === 0) {
      throw new Error("invalid priorities passed");
    }

    if (uniqueValidValues.length === 1) {
      return uniqueValidValues[0];
    }

    return uniqueValidValues;
  };
  return z.union([
    z.string().refine((value) => validValues.includes(value)),
    z.array(z.string()).transform(sanitizeArray),
    z.undefined(),
  ]);
};

const assigneesScheme = z.union([
  z.string(),
  z.array(z.string()),
  z.undefined(),
]);

// TODO: deprecate this usage, as we have modified the search params validations. 
export const defaultSortOp = "last-message-dsc";

// const sortEnum = z.enum(["last-message-dsc", "created-asc", "created-dsc"]);
// const sortSchema = z.union([sortEnum, z.undefined()]);

// using fallback to avoid the unknown
// see https://tanstack.com/router/latest/docs/framework/react/guide/search-params#:~:text=However%20the%20use%20of%20catch%20here%20overrides%20the%20types%20and%20makes%20page%2C%20filter%20and%20sort%20unknown%20causing%20type%20loss.%20We%20have%20handled
export const threadSearchSchema = z.object({
  stages: fallback(
    stagesSchema([
      "needs_first_response",
      "waiting_on_customer",
      "hold",
      "needs_next_response",
    ]),
    undefined
  ).catch(undefined),
  sort: fallback(
    z.enum(["last-message-dsc", "created-asc", "created-dsc"]),
    defaultSortOp
  ).catch(defaultSortOp),
  priorities: fallback(
    prioritiesSchema(["urgent", "high", "normal", "low"]),
    undefined
  ).catch(undefined),
  assignees: fallback(assigneesScheme, undefined).catch(undefined),
});

export type ThreadSearch = z.infer<typeof threadSearchSchema>;
