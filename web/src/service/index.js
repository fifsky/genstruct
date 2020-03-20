import {createApi} from "../util";

export const genApi = data => createApi("/genapi/struct/gen", data);
