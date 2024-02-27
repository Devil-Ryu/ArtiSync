// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {models} from '../models';
import {context} from '../models';

export function AutoMigrate():Promise<void>;

export function CheckConnect():Promise<void>;

export function ClearInterfaceRecord():Promise<void>;

export function Connect(arg1:string):Promise<void>;

export function CreateInterface(arg1:models.Interface):Promise<models.Interface>;

export function CreateOrUpdateInterfaceRecord(arg1:models.InterfaceRecord):Promise<models.InterfaceRecord>;

export function CreatePlatform(arg1:Array<models.Platform>):Promise<Array<models.Platform>>;

export function DeleteInterfaceRecord(arg1:models.InterfaceRecord):Promise<void>;

export function DeleteInterfaces(arg1:Array<models.Interface>):Promise<void>;

export function DeletePlatforms(arg1:Array<models.Platform>):Promise<void>;

export function ExportPlatform(arg1:string,arg2:models.Platform):Promise<void>;

export function GetConfigDir():Promise<string>;

export function GetConfigFilePath():Promise<string>;

export function GetInterface(arg1:models.Interface):Promise<models.Interface>;

export function GetInterfaceRecords(arg1:{[key: string]: any}):Promise<Array<models.InterfaceRecord>>;

export function GetPlatform(arg1:models.Platform):Promise<models.Platform>;

export function GetPlatforms(arg1:{[key: string]: any}):Promise<Array<models.Platform>>;

export function ImportPlatform(arg1:string):Promise<void>;

export function LoadJSONFile(arg1:string):Promise<{[key: string]: any}>;

export function QueryInterfaceRecords(arg1:{[key: string]: any},arg2:number,arg3:number):Promise<{[key: string]: any}>;

export function SaveJSONFile(arg1:string,arg2:{[key: string]: any}):Promise<void>;

export function Startup(arg1:context.Context):Promise<void>;

export function UpdateInterface(arg1:Array<models.Interface>):Promise<void>;

export function UpdatePlatform(arg1:models.Platform):Promise<void>;
