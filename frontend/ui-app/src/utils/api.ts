export type ApiResult<T, E> = {
    ok: T;
    err: E;
};

export const get = async <T, E>(url: string) => {
    const res = await fetch(url);
    if (!res.ok) {
        return (await res.json()) as ApiResult<T, E>;
    }
    return (await res.json()) as ApiResult<T, E>;
};
