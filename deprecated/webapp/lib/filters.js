class QueryFilter {
  constructor(params) {
    this.params = params;
    this.redirect = false;
    this.updatedParams = {};
  }

  static singleKV = ["status"];
  static statuses = ["todo", "snoozed", "done", "unsnoozed"];
  static defaultStatus = "todo";

  static reasons = ["replied", "unreplied"];
  static priorities = [0, 1, 2, 3]; // [urgent, high, normal, low]

  buildStatusQuery() {
    const { status } = this.params;
    if (!status || !QueryFilter.statuses.includes(status)) {
      this.updatedParams.status = QueryFilter.defaultStatus;
      this.redirect = true;
    }
    return this;
  }

  buildReasonsQuery() {
    const { reason = [] } = this.params;
    const uniqueReasons = new Set(reason);
    const t = [];
    let seen = false;
    for (const r of QueryFilter.reasons) {
      if (uniqueReasons.has(r)) {
        t.push(r);
        seen = true;
      }
    }
    if (seen) {
      this.updatedParams.reason = t;
    }
    return this;
  }

  buildPriorityQuery() {
    const { priority = [] } = this.params;
    const uniquePriorities = new Set(priority);
    const t = [];
    let seen = false;
    for (const r of QueryFilter.priorities) {
      if (uniquePriorities.has(r)) {
        t.push(r);
        seen = true;
      }
    }
    if (seen) {
      this.updatedParams.priority = t;
    }
    return this;
  }

  mergeParams() {
    this.updatedParams = { ...this.params, ...this.updatedParams };
    return this;
  }

  buildQuery() {
    this.buildStatusQuery().buildReasonsQuery().mergeParams();
    return this.generateQueryParams();
  }

  generateQueryParams() {
    const params = new URLSearchParams();
    for (const key in this.updatedParams) {
      const value = this.updatedParams[key];
      // check if its an Array
      if (Array.isArray(value)) {
        for (const v of value) {
          params.append(key, v);
        }
        continue;
      }
      params.append(key, value);
    }
    return params;
  }

  buildCleanedQuery() {
    this.buildStatusQuery().buildReasonsQuery();
    return this.generateQueryParams();
  }

  getStatus() {
    return this.updatedParams.status || QueryFilter.defaultStatus;
  }

  getReasons() {
    return this.updatedParams?.reason;
  }

  removeReason(reason) {
    const appliedReasons = this.updatedParams?.reason;
    if (!appliedReasons || !appliedReasons.length) {
      delete this.updatedParams.reason;
      return this;
    }

    const reasonsFiltered = appliedReasons.filter((r) => r !== reason);
    if (reasonsFiltered.length === 0) {
      delete this.updatedParams.reason;
      return this;
    }

    this.updatedParams.reason = reasonsFiltered;
    return this;
  }

  addReason(r) {
    const { reason = [] } = this.updatedParams;
    if (QueryFilter.reasons.includes(r)) {
      const reasons = new Set([...reason, r]);
      this.updatedParams.reason = Array.from(reasons);
    }
    return this;
  }

  static fromClientQueryParams(searchParams) {
    const queryParams = {};
    for (const key of searchParams.keys()) {
      if (QueryFilter.singleKV.includes(key)) {
        queryParams[key] = searchParams.get(key);
      } else {
        const values = searchParams.getAll(key);
        queryParams[key] = values;
      }
    }
    return queryParams;
  }
}

export { QueryFilter };
