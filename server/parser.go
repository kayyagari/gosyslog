package server

import (
   //"strconv"
   sysmsg "gosyslog/message"
)


func Parse(buf []byte) (sysmsg.Message, error){
  msg := sysmsg.Message{}
  return msg, nil
}