import * as grpc from '@grpc/grpc-js'
import { loadSync } from "@grpc/proto-loader";


const PROTO_PATH = './auth.proto';
const packageDefinition = loadSync(PROTO_PATH);
const authProto = grpc.loadPackageDefinition(packageDefinition).auth as any;
const authServiceClient = new authProto.AuthService(
    'localhost:50051',
    grpc.credentials.createInsecure()
);


export async function AuthenticateUser(req : any, res : any, next : any){
    const token = req.headers['authorization']?.replace('Bearer ', '').trim();
    if (!token) {
        return res.status(401).json({ error: 'Unauthorized: Missing token' });
    }
    
    authServiceClient.ValidateTokenUser({ token }, (err: any, response : any) => {
        if (err || !response.isValid) {
            return res.status(401).json({ error: response?.error || 'Unauthorized' });
        }

        req.user = {
            id: Number(response.userId),
        };

        next();
    });
}
export async function AuthenticateAdmin(req : any, res : any, next : any){
    const token = req.headers['authorization']?.replace('Bearer ', '').trim();
    if (!token) {
        return res.status(401).json({ error: 'Unauthorized: Missing token' });
    }
    
    authServiceClient.ValidateTokenAdmin({ token }, (err: any, response : any) => {
        if (err || !response.isValid) {
            return res.status(401).json({ error: response?.error || 'Unauthorized' });
        }

        req.admin = {
            id: Number(response.adminId),
        };

        next();
    });
}
