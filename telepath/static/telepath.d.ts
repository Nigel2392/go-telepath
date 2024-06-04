declare class Telepath {
    constructors: {
        [key: string]: any;
    };
    constructor();
    register(name: any, constructor: any): void;
    unpack(objData: any): any;
    scanForIds(objData: any, packedValuesById: {
        [key: number]: any;
    }): void;
    unpackWithRefs(objData: any, packedValuesById: {
        [key: number]: any;
    }, valuesById: {
        [key: number]: any;
    }): any;
}
export { Telepath, };
