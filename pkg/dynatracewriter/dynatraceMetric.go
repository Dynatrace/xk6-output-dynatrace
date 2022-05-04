package dynatracewriter

import (
   "time"
   "fmt"
   "strcomv"
   "go.k6.io/k6/output"
    "go.k6.io/k6/stats"
)

const {
    metricDisplayNameProperty="dt.meta.displayName"
    metricDescriptionProperty="dt.meta.description"
    metricUnitProperty="dt.meta.unit"
}
type dynatraceMetric struct{
    metricDisplayName string
    description string
    metricKeyName string
    metricUnit string
    metricDimensions map[string]string
    metricValue float64
    metricTimeStamp uint32
}

func New(metricDisplayName string, description string, metricKeyName string, metricUnit string,metricDimensions map[string]string,value float64,metricTimeStamp uint32) dynatraceMetric {
    e := dynatraceMetric {metricDisplayName, description, metricKeyName, metricUnit,metricDimensions,value,metricTimeStamp}
    return e
}
func samleToDynametric(sample Sample ) dynatraceMetric {
    e := dynatraceMetric { metricKeyName:= sample.Metric.Name, metricDimensions:=sample.CloneTags(), value:= sample.Value , metricTimeStamp:=sample.Time.UnixMilli()   }
    return e
}


func toText(e dynatraceMetric) string{

   result string

   result=e.metricKeyName

   if(len(e.metricDimensions)!=0) {
        for key, value := range e.metricDimensions {
                result+=","+key+ "="+ value
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
    result+=" "+strconv.Itoa(e.metricTimeStamp)

    return result
}