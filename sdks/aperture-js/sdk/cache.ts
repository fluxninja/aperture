export type CachedValueResponse = {
    error: Error | null;
    lookupResult: string | null;
    code: string | null;
    message: string | null;
    value: Buffer | null;
};

export type SetCachedValueResponse = {
    error: Error | null;
    code: string | null;
    message: string | null;
};

export type DeleteCachedValueResponse = {
    error: Error | null;
    code: string | null;
    message: string | null;
};
