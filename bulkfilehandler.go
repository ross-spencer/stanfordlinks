package main 

import (
  	"os"
   "fmt"
)

type contenterror struct {
   content string
   err error
}  

func extractAndAnalyse(filepool []filedata) (bool, error) {
   
   ch := make(chan contenterror)
   for _, fi := range filepool {
      go getFileContent(fi, ch)
   }

   for _, fi := range filepool{
      ce := <- ch
      if ce.err != nil {
         logFileMessage("INFO: '%s' cannot be handled by Tika.", fi.fname)
      } else {
         edat := getEntityData(ce.content, fi.fname) 
         go collateEntities(edat)         
      }
   }
   return false, nil
}

func openFile (path string) (*os.File, error) {
   fp, err := os.Open(path)
   if err != nil {
      return nil, err
   }
   return fp, nil
}

func getFileContent(fi filedata, ch chan contenterror) {
   var ce contenterror

   //what are we doing..?
   logFileMessage("INFO: '%s' being processed.", fi.fname)

   //process...
   fp, err := openFile(fi.fpath)
   defer fp.Close()
   if err != nil {
      ce.err = err
      ch <- ce
      return
   }

   _, fl_recursive_keys_values, err := getTikaRecursive(fi.fname, fp, ACCEPT_MIME_JSON)
   if err != nil {
      ce.err = err
      ch <- ce
      return
   }

   if val, ok := fl_recursive_keys_values[TIKA_PLAIN_TEXT]; ok {
      ce.content = val.(string)
      ch <- ce
      return
   } else {
      ce.err = fmt.Errorf("No plain text data to analyse via NER")
      ch <- ce
      return 
   }
}

func getEntityData(content string, fname string) []EntityData {
   return getNERData(content, fname)
}

