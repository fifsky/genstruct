import {createApi} from "../util";

export const genApi = data => createApi("/api/struct/gen", data);
