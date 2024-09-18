import { z } from "zod";

//
// for more: https://tanstack.com/router/latest/docs/framework/react/guide/search-params
// usage of `.catch` or `default` matters.
const reasonsSchema = (validValues: string[]) => {
  const sanitizeArray = (arr: string[]) => {
    // remove duplicates
    const uniqueValues = [...new Set(arr)];
    // filter only valid values
    const uniqueValidValues: string[] = uniqueValues.filter((val) =>
      validValues.includes(val)
    );

    // no valid values
    if (uniqueValidValues.length === 0) {
      throw new Error("invalid reason(s) passed");
    }

    if (uniqueValidValues.length === 1) {
      return uniqueValidValues[0];
    }

    return uniqueValidValues;
  };
  return z.union([
    z.string().refine((value) => validValues.includes(value)),
    z.array(z.string()).transform(sanitizeArray),
    // .refine(
    //   (arr) =>
    //     arr.length === validValues.length &&
    //     validValues.every((val) => arr.includes(val))
    // ),
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

// Defaults
export const defaultSortOp = "last-message-dsc";

// const sortEnum = z.enum(["last-message-dsc", "created-asc", "created-dsc"]);
// const sortSchema = z.union([sortEnum, z.undefined()]);

export const threadSearchSchema = z.object({
  reasons: reasonsSchema(["replied", "unreplied"]).catch(""),
  // sort: sortEnum.catch(defaultSortOp),
  sort: z
    .enum(["last-message-dsc", "created-asc", "created-dsc"])
    .catch("last-message-dsc"),
  priorities: prioritiesSchema(["urgent", "high", "normal", "low"]).catch(""),
  assignees: assigneesScheme.catch(""),
});
