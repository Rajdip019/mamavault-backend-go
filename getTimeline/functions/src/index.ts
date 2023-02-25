import { getTime } from "date-fns";
import * as functions from "firebase-functions";

interface IDocuments{
    doc_id : string,
    doc_title : string,
    doc_type : string,
    doc_format : string,
    doc_url : string,
    doc_download_url : string
    timeline_time : string,
    upload_time : string
}

export const getTimeline = functions.https.onRequest((request, response) => {
    const userDocuments : IDocuments[]  = request.body.documents
    const date_filterd_docs :  {time : number, document : IDocuments[]}[]= []
    const dates = userDocuments.map((data : IDocuments) => { return getTime(new Date(data.timeline_time)) })
    dates.map((date : number) => {
        userDocuments.map((data :IDocuments )=> {
            if (date === getTime(new Date(data.timeline_time))){
                    if(date_filterd_docs.filter(e => getTime(new Date(data.timeline_time)) === e.time ).length > 0){
                        const index = date_filterd_docs.findIndex(e => getTime(new Date(data.timeline_time)) === e.time)
                        if(date_filterd_docs[index].document.includes(data)){
                            console.log("Already contains data");
                        }else{
                            date_filterd_docs[index].document.push(data)
                        }
                    }else{
                        date_filterd_docs.push({
                            time : date,
                            document : [data]
                        })
                    }
            }
        }) 
    })
    date_filterd_docs.sort((a, b) => {
        return b.time - a.time
    })
    response.send(date_filterd_docs);
});
