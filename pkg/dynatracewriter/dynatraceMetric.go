package dynatracewriter

import (
   "time"
   "fmt"
   "strconv"
    "go.k6.io/k6/metrics"
)

const (
    metricDisplayNameProperty="dt.meta.displayName"
    metricDescriptionProperty="dt.meta.description"
    metricUnitProperty="dt.meta.unit"
    metricKeyPrefix="k6"
)
type dynatraceMetric struct{
    metricDisplayName string
    description string
    metricKeyName string
    metricUnit string
    metricDimensions map[string]string
    metricValue float64
    metricTimeStamp int64
}


func samleToDynametric(sample metrics.Sample ) dynatraceMetric {
     return dynatraceMetric{
        metricKeyName : sample.Metric.Name,
        metricDimensions : sample.GetTags().Map(),
        metricValue : sample.Value,
        metricTimeStamp : sample.GetTime().UnixMilli(),
     }
}


func (e *dynatraceMetric) toText() string {

   var result=""

   result=metricKeyPrefix+"."+e.metricKeyName

   if(len(e.metricDimensions)!=0) {
        for key, value := range e.metricDimensions {
                if len(key)>0 && len(value)>0 {
                     result+=","+key+ "="+ "\""+value+"\""
                }
        }
   }

    result+=" "
   if(len(e.metricUnit)>0){
        result+=metricUnitProperty+"="+e.metricUnit
    }

    if(len(e.description)>0){
        result+=","+metricDescriptionProperty+"="+e.description

    }

    if(len(e.metricDisplayName)>0){
            result+=","+metricDisplayNameProperty+"="+e.metricDisplayName
    }

    result+=" "+ fmt.Sprint(e.metricValue)

    if e.metricTimeStamp<= 0 {
        t := time.Now() //It will return time.Time object with current timestamp
        tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
        e.metricTimeStamp=tUnixMilli

    }
    result+=" "+strconv.FormatInt(e.metricTimeStamp,10)

    return result
}
