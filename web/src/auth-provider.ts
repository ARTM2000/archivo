import { AuthProvider as IAuthProvider } from "ra-core";

export const AuthProvider: IAuthProvider = {
    login: async (params) => {
        console.log(params);
        
    },
    checkAuth: async (params) => {
        throw new Error()
    },
    checkError: async (error) => {
        
    },
    logout: async () => {
        
    },
    getIdentity: async () => {
        return {
            id: 0,
        }
    },
    getPermissions: async () => {},
}
