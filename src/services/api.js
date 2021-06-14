class BaseApiService{
    baseUrl = process.env.VUE_APP_API_URL_LOAD_DATA_BY_DATE;
    resource;


    constructor(resource){
        if (!resource) throw new Error("Resource not provided");
        this.resource = resource;
    }
    
}