"use client";

import { QueryFilter } from "@/lib/filters";
import * as React from "react";

import { useRouter, useSearchParams } from "next/navigation";

function useReasonFilter() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [selectedReasons, setSelectedReasons] = React.useState([]);

  function setReason(reason) {
    const queryParams = QueryFilter.fromClientQueryParams(searchParams);
    const filters = new QueryFilter(queryParams);
    const paramsUpdated = filters
      .buildReasonsQuery()
      .mergeParams()
      .addReason(reason)
      .generateQueryParams();

    router.push(`?${paramsUpdated.toString()}`, undefined, {
      shallow: true,
    });
  }

  function clearReason(reason) {
    const queryParams = QueryFilter.fromClientQueryParams(searchParams);
    const filters = new QueryFilter(queryParams);
    const paramsUpdated = filters
      .buildReasonsQuery()
      .mergeParams()
      .removeReason(reason)
      .generateQueryParams();

    router.push(`?${paramsUpdated.toString()}`, undefined, {
      shallow: true,
    });
  }

  function isChecked(reason) {
    return selectedReasons.includes(reason);
  }

  function countSelectedReasons() {
    return selectedReasons.length;
  }

  React.useEffect(() => {
    const queryParams = QueryFilter.fromClientQueryParams(searchParams);
    const filters = new QueryFilter(queryParams);
    const reasons = filters.buildReasonsQuery().getReasons();
    if (reasons) setSelectedReasons([...reasons]);
    else setSelectedReasons([]);
  }, [searchParams]);

  return {
    setReason,
    clearReason,
    isChecked,
    countSelectedReasons,
  };
}

export { useReasonFilter };
