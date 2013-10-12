

/*2:*/


//line amail.w:32

//line license:1

// This file is part of Amail version 0.4
// Author Alexander Sychev
//
// Copyright (c) 2013 Alexander Sychev. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * The name of author may not be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main
//line amail.w:34
 import(


/*11:*/


//line amail.w:129

"flag"
"fmt"
"os"
"strings"
"sort"



/*:11*/



/*13:*/


//line amail.w:163

"unicode"
"unicode/utf8"



/*:13*/



/*15:*/


//line amail.w:181

"code.google.com/p/goplan9/plan9/client"
"github.com/golang/glog"



/*:15*/



/*28:*/


//line amail.w:309

"io"
"bufio"



/*:28*/



/*31:*/


//line amail.w:355

"bitbucket.org/santucco/goplumb"
"code.google.com/p/goplan9/plan9"



/*:31*/



/*38:*/


//line amail.w:501

"bitbucket.org/santucco/goacme"



/*:38*/



/*47:*/


//line amail.w:652

"strconv"



/*:47*/



/*91:*/


//line amail.w:1221

"time"



/*:91*/



/*200:*/


//line amail.w:2634

"os/exec"



/*:200*/



/*215:*/


//line amail.w:2965

"os/user"



/*:215*/



/*219:*/


//line amail.w:2986

"sync"



/*:219*/


//line amail.w:35

)

type(


/*19:*/


//line amail.w:205

mailbox struct{
name string


/*21:*/


//line amail.w:231

all messages
unread messages
mch chan int
dch chan int



/*:21*/



/*48:*/


//line amail.w:656

fid*client.Fid
total int



/*:48*/



/*72:*/


//line amail.w:948

shownew bool
showthreads bool
ech<-chan*goacme.Event
w*goacme.Window
cch chan bool



/*:72*/



/*77:*/


//line amail.w:999

thread bool



/*:77*/



/*101:*/


//line amail.w:1316

markch chan[]int
unmarkch chan[]int



/*:101*/



/*107:*/


//line amail.w:1406

mdch chan messages



/*:107*/



/*136:*/


//line amail.w:1709

rfch chan*refresh
irfch chan*refresh
reset bool



/*:136*/



/*152:*/


//line amail.w:1869

pos int



/*:152*/



/*166:*/


//line amail.w:2040

deleted int



/*:166*/



/*178:*/


//line amail.w:2247

lch chan[]int



/*:178*/


//line amail.w:208

}

mailboxes[]*mailbox

message struct{
id int


/*29:*/


//line amail.w:314

unread bool
box*mailbox



/*:29*/



/*64:*/


//line amail.w:880

deleted bool



/*:64*/



/*92:*/


//line amail.w:1225

from string
date time.Time
subject string



/*:92*/



/*116:*/


//line amail.w:1466

inreplyto string
messageid string



/*:116*/



/*119:*/


//line amail.w:1503

parent*message



/*:119*/



/*180:*/


//line amail.w:2255

w*goacme.Window



/*:180*/



/*187:*/


//line amail.w:2341

to[]string
cc[]string



/*:187*/



/*202:*/


//line amail.w:2656

text string
html string
showhtml bool
files[]*file
cids map[string]*file



/*:202*/


//line amail.w:215

}

messages[]*message




/*:19*/



/*82:*/


//line amail.w:1149

msgmap map[string][]int



/*:82*/



/*114:*/


//line amail.w:1457

idmessages[]*message



/*:114*/



/*134:*/


//line amail.w:1690

refreshFlags int

refresh struct{
flags refreshFlags
msgs messages
}



/*:134*/



/*201:*/


//line amail.w:2638

file struct{
name string
mimetype string
path string
}



/*:201*/


//line amail.w:39

)



/*41:*/


//line amail.w:530

const mailboxfmt= "%-30s\t%10d\t%10d\n"
const wholefile= "0,$"



/*:41*/



/*45:*/


//line amail.w:612

const mailboxfmtprc= "%-30s\t%10d\t%10d\t%d%%\n"



/*:45*/



/*135:*/


//line amail.w:1701

const(
seek refreshFlags= 1<<iota
insert refreshFlags= 1<<iota
exact refreshFlags= 1<<iota
)



/*:135*/



/*156:*/


//line amail.w:1903

const eof= "$"



/*:156*/



/*171:*/


//line amail.w:2101

const bof= "#0-"
const eol= "+#0"



/*:171*/


//line amail.w:42


var(


/*3:*/


//line amail.w:78

exit chan bool= make(chan bool)



/*:3*/



/*6:*/


//line amail.w:96

wch chan int= make(chan int,100)
wcount int



/*:6*/



/*10:*/


//line amail.w:121

shownew bool
showthreads bool
levelmark string
newmark string
skipboxes[]string



/*:10*/



/*16:*/


//line amail.w:186

fsys*client.Fsys
rfid*client.Fid
srv string= "mail"



/*:16*/



/*20:*/


//line amail.w:222

boxes mailboxes



/*:20*/



/*23:*/


//line amail.w:249

mch= make(chan*struct{name string;id int},100)
dch= make(chan*struct{name string;id int},100)
bch= make(chan string,10)
rfch= make(chan*mailbox,100)



/*:23*/



/*39:*/


//line amail.w:505

mw*goacme.Window
ech<-chan*goacme.Event



/*:39*/



/*44:*/


//line amail.w:608

shown= make(map[string]int)



/*:44*/



/*83:*/


//line amail.w:1153

lch= make(chan*msgmap,100)



/*:83*/



/*89:*/


//line amail.w:1204

deleted= "(deleted)-"



/*:89*/



/*97:*/


//line amail.w:1280

markch= make(chan*msgmap,100)
unmarkch= make(chan*msgmap,100)



/*:97*/



/*110:*/


//line amail.w:1432

mdch chan messages= make(chan messages,100)



/*:110*/



/*115:*/


//line amail.w:1461

idmap= make(map[string]*struct{msg*message;children idmessages})
idch= make(chan struct{id string;val interface{}})



/*:115*/



/*150:*/


//line amail.w:1855

mrfch chan*refresh= make(chan*refresh)



/*:150*/



/*207:*/


//line amail.w:2802

home string



/*:207*/



/*216:*/


//line amail.w:2969

cuser string



/*:216*/



/*220:*/


//line amail.w:2990

once sync.Once



/*:220*/



/*230:*/


//line amail.w:3275

plan9dir string



/*:230*/


//line amail.w:45

debug glog.Level= 1
)

func main(){
glog.V(debug).Infoln("main")
defer glog.V(debug).Infoln("main is done")


/*12:*/


//line amail.w:137

{
glog.V(debug).Infoln("parsing command line arguments")
var skip string
flag.BoolVar(&shownew,"new",true,"show new messages only")
flag.BoolVar(&showthreads,"threads",true,"show threads of messages")
flag.StringVar(&skip,"skip","","boxes to be skiped, separated by comma")
flag.StringVar(&levelmark,"levelmark","+","mark of level for threads")
flag.StringVar(&newmark,"newmark","(*)","mark of new messages")
flag.Usage= func(){
fmt.Fprintf(os.Stderr,"Mail client for Acme programming environment\n")
fmt.Fprintf(os.Stderr,"Usage: %s [options] [<mailbox 1>]...[<mailbox N>]\n",os.Args[0])
fmt.Fprintln(os.Stderr,"Options:")
flag.PrintDefaults()
}
flag.Parse()


/*14:*/


//line amail.w:169

glog.V(debug).Infoln("checking of levelmark and newmark")
if r,_:=utf8.DecodeLastRuneInString(levelmark);unicode.IsDigit(r){
fmt.Fprintln(os.Stderr,"last symbol of level mark shouldn't be a digit")
os.Exit(1)
}
if r,_:=utf8.DecodeRuneInString(newmark);unicode.IsDigit(r){
fmt.Fprintln(os.Stderr,"first symbol of new mark shouldn't be a digit")
os.Exit(1)
}



/*:14*/


//line amail.w:153

if len(skip)> 0{
skipboxes= strings.Split(skip,", ")
sort.Strings(skipboxes)
glog.V(debug).Infof("these mailboxes will be skipped: %v\n",skipboxes)

}
}



/*:12*/


//line amail.w:52



/*17:*/


//line amail.w:192

{
glog.V(debug).Infoln("try to open mailfs")
var err error
if fsys,err= client.MountService(srv);err!=nil{
glog.Errorf("can't mount mailfs: %v\n",err)
os.Exit(1)
}
}




/*:17*/


//line amail.w:53



/*32:*/


//line amail.w:362

{
glog.V(debug).Infoln("trying to open 'seemail' plumbing port")
if sm,err:=goplumb.Open("seemail",plan9.OREAD);err!=nil{
glog.Errorf("can't open plumb/seemail: %s\n",err)
}else{
sch,err:=sm.MessageChannel(0)
if err!=nil{
glog.Errorf("can't get message channal for plumb/seemail: %s\n",err)
}else{
go func(){
defer sm.Close()
defer glog.V(debug).Infoln("plumbing goroutine is done")
for{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:377

case m,ok:=<-sch:
if!ok{
glog.Warningln("it seems plumber has finished")
sch= nil
return
}
glog.V(debug).Infof("a plumbing message has been received: %v\n",m)
if m.Attr["filetype"]!="mail"{
glog.Warningln("attribute 'filetype' is not 'mail'")
continue
}
v,ok:=m.Attr["mailtype"]
if!ok{
glog.Warningln("can't find 'mailtype' attribute")
continue
}
b:=strings.Split(string(m.Data),"/")
if len(b)<3{
glog.Warningln("can't read a name of mailbox and a number of message")
continue
}
num,err:=strconv.Atoi(b[2])
if err!=nil{
glog.Error(err)
continue
}
if v=="new"{
glog.V(debug).Infof("'%d' is a new message in the '%s' mailbox\n",num,b[1])
mch<-&struct{name string;id int}{name:b[1],id:num}
}else if v=="delete"{
glog.V(debug).Infof("'%d' is a deleted message in the '%s' mailbox\n",num,b[1])
dch<-&struct{name string;id int}{name:b[1],id:num}
}
}
}
}()
}
}
}



/*:32*/


//line amail.w:54



/*25:*/


//line amail.w:273

glog.V(debug).Infoln("initialization of root of mailfs")
var err error
rfid,err= fsys.Walk(".")
if err!=nil{
glog.Errorf("can't open mailfs: %v\n",err)
os.Exit(1)
}
defer rfid.Close()




/*:25*/


//line amail.w:55



/*118:*/


//line amail.w:1485

go func(){
for{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:1489

case v:=<-idch:
if v.val==nil{


/*128:*/


//line amail.w:1605

{
val,ok:=idmap[v.id]
if!ok{
continue
}
for i,v:=range val.children{
glog.V(debug).Infof("clear the parent of the '%d'\n",v.id)
val.children[i].parent= nil
}
if val.msg!=nil&&val.msg.parent!=nil{
if p,ok:=idmap[val.msg.parent.messageid];ok{
for i,_:=range p.children{
if p.children[i]==val.msg{
glog.V(debug).Infof("removing the '%d' message from the children of the message '%d'\n",
val.msg.id,val.msg.parent.id)
p.children.Delete(i)
break
}
}
}
val.msg= nil
}
if len(val.children)==0{
delete(idmap,v.id)
}
}



/*:128*/


//line amail.w:1492

}else if msg,ok:=v.val.(*message);ok{


/*127:*/


//line amail.w:1570

{
glog.V(debug).Infof("appending a '%s' message to idmap\n",v.id)
if val,ok:=idmap[v.id];!ok{
glog.V(debug).Infof("'%s' message  doesn't exist, creating\n",v.id)
idmap[v.id]= &struct{msg*message;children idmessages}{msg,nil}
}else{
if val.msg!=nil{
continue
}
glog.V(debug).Infof("'%s' message exists, reseting\n",v.id)
val.msg= msg
for i,_:=range val.children{
glog.V(debug).Infof("setting '%s' like a parent of '%s'\n",msg.messageid,val.children[i].messageid)
val.children[i].parent= msg
}
}
if len(msg.inreplyto)==0{
continue
}
if val,ok:=idmap[msg.inreplyto];!ok{
glog.V(debug).Infof("'%s' message (parent of '%s') doesn't exist, creating\n",msg.inreplyto,msg.messageid)
idmap[msg.inreplyto]= &struct{msg*message;children idmessages}{nil,append(idmessages{},msg)}
}else{
glog.V(debug).Infof("'%s' message exists, appending the '%s' like a child\n",msg.inreplyto,msg.messageid)
if _,ok:=val.children.SearchInsert(msg);ok&&val.msg!=nil{
glog.V(debug).Infof("setting '%s' like a parent of '%s'\n",val.msg.messageid,msg.messageid)
msg.parent= val.msg
}
}
}



/*:127*/


//line amail.w:1494

}else if ch,ok:=v.val.(chan idmessages);ok{


/*130:*/


//line amail.w:1651

{
if val,ok:=idmap[v.id];ok{
glog.V(debug).Infof("sending children for '%s'\n",v.id)
children:=make(idmessages,len(val.children),len(val.children))
copy(children,val.children)
sort.Sort(children)
ch<-children
}else{
glog.V(debug).Infof("'%s' is not found\n",v.id)
ch<-nil
}
}



/*:130*/


//line amail.w:1496

}
}
}
}()



/*:118*/


//line amail.w:56

if len(flag.Args())> 0{


/*34:*/


//line amail.w:452

go func(){
glog.V(debug).Infoln("start a main message loop")
defer glog.V(debug).Infoln("main message loop is done")
for{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:458

case b:=<-rfch:
if b==nil{


/*42:*/


//line amail.w:535

if mw!=nil{
glog.V(debug).Infoln("printing of the mailboxes")
if err:=mw.WriteAddr(wholefile);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",wholefile,err)
}else if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else{
for _,v:=range boxes{
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,len(v.unread),len(v.all))))
}
}
w:=mw


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:548

}




/*:42*/


//line amail.w:461

}else{


/*46:*/


//line amail.w:617

glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n",b.name,len(b.all),b.total)
if mw!=nil{
if len(b.all)!=b.total{
if c,ok:=shown[b.name];!ok||c<99{
shown[b.name]= c+1
continue
}else{
shown[b.name]= 0
}
}

if err:=mw.WriteAddr("0/^%s.*\\n/",escape(b.name));err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
continue
}

if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if len(b.all)!=b.total{
if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,b.name,len(b.unread),len(b.all),len(b.all)*100/b.total)));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
}else if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmt,b.name,len(b.unread),len(b.all))));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
w:=mw


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:648

}



/*:46*/


//line amail.w:463

}


/*43:*/


//line amail.w:554

case ev,ok:=<-ech:
glog.V(debug).Infof("an event from main window has been received: %v\n",ev)
if!ok{
ech= nil


/*9:*/


//line amail.w:115

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:559

return
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"ShowNew":
shownew= true
case"ShowAll":
shownew= false
case"ShowThreads":
showthreads= true
case"ShowPlain":
showthreads= false
case"Del":
mw.UnreadEvent(ev)
mw.Close()
mw= nil


/*5:*/


//line amail.w:88

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:576

return
case"debug":
debug= 0
continue
case"nodebug":
debug= 1
continue
default:
mw.UnreadEvent(ev)
continue
}


/*164:*/


//line amail.w:2013

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:164*/


//line amail.w:588

continue
}else if(ev.Type&goacme.Look)==goacme.Look{
name:=ev.Text
// a box name can contain spaces
if len(ev.Arg)> 0{
name+= " "+ev.Arg
}
name= strings.TrimSpace(name)
if i,ok:=boxes.Search(name);ok{
box:=boxes[i]


/*74:*/


//line amail.w:962

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:74*/


//line amail.w:599

continue
}
}
mw.UnreadEvent(ev)




/*:43*/


//line amail.w:465

}
}
}()



/*:34*/


//line amail.w:58

for _,name:=range flag.Args(){


/*35:*/


//line amail.w:471

glog.V(debug).Infof("creating a '%s' mailbox\n",name)
box:=&mailbox{name:name,

/*22:*/


//line amail.w:238

mch:make(chan int,100),
dch:make(chan int,100),




/*:22*/



/*73:*/


//line amail.w:956

shownew:shownew,
showthreads:showthreads,
cch:make(chan bool,100),



/*:73*/



/*102:*/


//line amail.w:1321

markch:make(chan[]int,100),
unmarkch:make(chan[]int,100),



/*:102*/



/*108:*/


//line amail.w:1410

mdch:make(chan messages,100),



/*:108*/



/*137:*/


//line amail.w:1715

rfch:make(chan*refresh,100),
irfch:make(chan*refresh,100),



/*:137*/



/*179:*/


//line amail.w:2251

lch:make(chan[]int,100),



/*:179*/


//line amail.w:473
}
boxes= append(boxes,box)
sort.Sort(boxes)



/*:35*/


//line amail.w:60



/*60:*/


//line amail.w:818

go box.loop()



/*:60*/


//line amail.w:61



/*74:*/


//line amail.w:962

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:74*/


//line amail.w:62

}
}else{


/*40:*/


//line amail.w:510

glog.V(debug).Infoln("creating the main window")
defer goacme.DeleteAll()

var err error
if mw,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
name:="Amail"
w:=mw


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:521

if ech,err= mw.EventChannel(0,goacme.Mouse,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*164:*/


//line amail.w:2013

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:164*/


//line amail.w:526



/*8:*/


//line amail.w:110

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:527




/*:40*/


//line amail.w:65



/*34:*/


//line amail.w:452

go func(){
glog.V(debug).Infoln("start a main message loop")
defer glog.V(debug).Infoln("main message loop is done")
for{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:458

case b:=<-rfch:
if b==nil{


/*42:*/


//line amail.w:535

if mw!=nil{
glog.V(debug).Infoln("printing of the mailboxes")
if err:=mw.WriteAddr(wholefile);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",wholefile,err)
}else if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else{
for _,v:=range boxes{
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,len(v.unread),len(v.all))))
}
}
w:=mw


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:548

}




/*:42*/


//line amail.w:461

}else{


/*46:*/


//line amail.w:617

glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n",b.name,len(b.all),b.total)
if mw!=nil{
if len(b.all)!=b.total{
if c,ok:=shown[b.name];!ok||c<99{
shown[b.name]= c+1
continue
}else{
shown[b.name]= 0
}
}

if err:=mw.WriteAddr("0/^%s.*\\n/",escape(b.name));err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
continue
}

if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if len(b.all)!=b.total{
if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,b.name,len(b.unread),len(b.all),len(b.all)*100/b.total)));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
}else if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmt,b.name,len(b.unread),len(b.all))));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
w:=mw


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:648

}



/*:46*/


//line amail.w:463

}


/*43:*/


//line amail.w:554

case ev,ok:=<-ech:
glog.V(debug).Infof("an event from main window has been received: %v\n",ev)
if!ok{
ech= nil


/*9:*/


//line amail.w:115

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:559

return
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"ShowNew":
shownew= true
case"ShowAll":
shownew= false
case"ShowThreads":
showthreads= true
case"ShowPlain":
showthreads= false
case"Del":
mw.UnreadEvent(ev)
mw.Close()
mw= nil


/*5:*/


//line amail.w:88

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:576

return
case"debug":
debug= 0
continue
case"nodebug":
debug= 1
continue
default:
mw.UnreadEvent(ev)
continue
}


/*164:*/


//line amail.w:2013

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:164*/


//line amail.w:588

continue
}else if(ev.Type&goacme.Look)==goacme.Look{
name:=ev.Text
// a box name can contain spaces
if len(ev.Arg)> 0{
name+= " "+ev.Arg
}
name= strings.TrimSpace(name)
if i,ok:=boxes.Search(name);ok{
box:=boxes[i]


/*74:*/


//line amail.w:962

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:74*/


//line amail.w:599

continue
}
}
mw.UnreadEvent(ev)




/*:43*/


//line amail.w:465

}
}
}()



/*:34*/


//line amail.w:66

go func(){


/*26:*/


//line amail.w:285

{
glog.V(debug).Infoln("enumerating of mailboxes")
fi,err:=rfid.Dirreadall()
if err!=nil{
glog.Errorf("can't read mailfs: %v\n",err)


/*5:*/


//line amail.w:88

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:291

return
}
for _,f:=range fi{
if f.Mode&plan9.DMDIR==plan9.DMDIR{
name:=f.Name


/*27:*/


//line amail.w:304

glog.V(debug).Infof("send a mailbox '%s' to put in the list\n",name)
bch<-name



/*:27*/


//line amail.w:297

}
}
glog.V(debug).Infoln("enumerating of mailboxes is done")
}



/*:26*/


//line amail.w:68

}()
}


/*33:*/


//line amail.w:423

glog.V(debug).Infoln("process events are specific for the list of mailboxes")
for{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:427

case name:=<-bch:


/*54:*/


//line amail.w:754

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:429



/*35:*/


//line amail.w:471

glog.V(debug).Infof("creating a '%s' mailbox\n",name)
box:=&mailbox{name:name,

/*22:*/


//line amail.w:238

mch:make(chan int,100),
dch:make(chan int,100),




/*:22*/



/*73:*/


//line amail.w:956

shownew:shownew,
showthreads:showthreads,
cch:make(chan bool,100),



/*:73*/



/*102:*/


//line amail.w:1321

markch:make(chan[]int,100),
unmarkch:make(chan[]int,100),



/*:102*/



/*108:*/


//line amail.w:1410

mdch:make(chan messages,100),



/*:108*/



/*137:*/


//line amail.w:1715

rfch:make(chan*refresh,100),
irfch:make(chan*refresh,100),



/*:137*/



/*179:*/


//line amail.w:2251

lch:make(chan[]int,100),



/*:179*/


//line amail.w:473
}
boxes= append(boxes,box)
sort.Sort(boxes)



/*:35*/


//line amail.w:430



/*71:*/


//line amail.w:936

glog.V(debug).Infoln("sending a signal to refresh all mailboxes")
rfch<-nil



/*:71*/


//line amail.w:431



/*60:*/


//line amail.w:818

go box.loop()



/*:60*/


//line amail.w:432

case d:=<-mch:
name:=d.name


/*37:*/


//line amail.w:490

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:754

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:494

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:435

glog.V(debug).Infof("sending '%d' to add in the '%s' mailbox\n",d.id,boxes[i].name)
boxes[i].mch<-d.id
case d:=<-dch:
name:=d.name


/*37:*/


//line amail.w:490

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:754

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:494

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:440

glog.V(debug).Infof("sending '%d' to delete from the '%s' mailbox\n",d.id,boxes[i].name)
boxes[i].dch<-d.id


/*7:*/


//line amail.w:101

case i:=<-wch:
wcount+= i
if wcount==0{


/*5:*/


//line amail.w:88

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:105

return
}



/*:7*/



/*86:*/


//line amail.w:1166

case d:=<-lch:
if d==nil{
continue
}
for name,ids:=range*d{


/*37:*/


//line amail.w:490

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:754

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:494

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:1172

boxes[i].lch<-ids
}



/*:86*/



/*98:*/


//line amail.w:1285

case d:=<-markch:
if d==nil{
continue
}
for name,ids:=range*d{


/*37:*/


//line amail.w:490

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:754

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:494

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:1291

boxes[i].markch<-ids
}
case d:=<-unmarkch:
if d==nil{
continue
}
for name,ids:=range*d{


/*37:*/


//line amail.w:490

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:754

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:494

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:1299

boxes[i].unmarkch<-ids
}




/*:98*/



/*111:*/


//line amail.w:1436

case msgs:=<-mdch:
for i,_:=range boxes{
glog.V(debug).Infof("sending %d messages to delete in the '%s' mailbox\n",len(msgs),boxes[i].name)
boxes[i].mdch<-append(messages{},msgs...)
}



/*:111*/



/*151:*/


//line amail.w:1859

case r:=<-mrfch:
for i,_:=range boxes{
glog.V(debug).Infof("sending messages to refresh in the '%s' mailbox\n",boxes[i].name)
boxes[i].rfch<-&refresh{r.flags,append(messages{},r.msgs...)}
}





/*:151*/


//line amail.w:443

}
}




/*:33*/


//line amail.w:71

}



/*:2*/



/*24:*/


//line amail.w:257

func(this mailboxes)Len()int{
return len(this)
}

func(this mailboxes)Less(i,j int)bool{
return this[i].name<this[j].name
}

func(this mailboxes)Swap(i,j int){
t:=this[i]
this[i]= this[j]
this[j]= t
}



/*:24*/



/*30:*/


//line amail.w:320

func(this*mailbox)newMessage(id int)(msg*message,unread bool,err error){
glog.V(debug).Infof("newMessage: trying to open '%d/info'\n",id)
f,err:=this.fid.Walk(fmt.Sprintf("%d/info",id))
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open to '%s/%d/info': %s\n",this.name,id,err)
return
}
defer f.Close()
msg= &message{id:id,box:this,

/*203:*/


//line amail.w:2664

cids:make(map[string]*file),



/*:203*/


//line amail.w:332
}
b:=bufio.NewReader(f)
unread= true
glog.V(debug).Infof("newMessage: reading and parsing of a content of '%d/info'\n",id)
for s,err:=b.ReadString('\n');err==nil;s,err= b.ReadString('\n'){
if s[len(s)-1]=='\n'{
s= s[:len(s)-1]
}
if strings.HasPrefix(s,"flags "){
if strings.Index(s,"seen")>=0{
unread= false
}
continue
}


/*93:*/


//line amail.w:1231

if strings.HasPrefix(s,"from "){
msg.from= s[len("from "):]
msg.from= strings.Replace(msg.from,"'' ","",-1)
continue
}
var unixdate int64
if _,err:=fmt.Sscanf(s,"unixdate %d",&unixdate);err==nil{
msg.date= time.Unix(unixdate,0)
continue
}
if strings.HasPrefix(s,"subject "){
msg.subject= s[len("subject "):]
continue
}




/*:93*/



/*117:*/


//line amail.w:1471

{
if _,err:=fmt.Sscanf(s,"inreplyto %s",&msg.inreplyto);err==nil{
msg.inreplyto= strings.Trim(msg.inreplyto,"<>")
continue
}
if _,err:=fmt.Sscanf(s,"messageid %s",&msg.messageid);err==nil{
msg.messageid= strings.Trim(msg.messageid,"<>")
idch<-struct{id string;val interface{}}{msg.messageid,msg}
continue
}
}



/*:117*/



/*188:*/


//line amail.w:2346

if strings.HasPrefix(s,"to "){
msg.to= split(s[len("to "):])
continue
}
if strings.HasPrefix(s,"cc "){
msg.cc= split(s[len("cc "):])
continue
}




/*:188*/


//line amail.w:346

}
msg.unread= unread
return

}



/*:30*/



/*36:*/


//line amail.w:479

func(this mailboxes)Search(name string)(int,bool){
pos:=sort.Search(len(this),
func(i int)bool{return this[i].name>=name});
if pos!=len(this)&&this[pos].name==name{
return pos,true
}
return pos,false
}



/*:36*/



/*52:*/


//line amail.w:735

func escape(s string)(res string){
for _,v:=range s{
if strings.ContainsRune("\\/[].+?()*^$",v){
res+= "\\"
}
res+= string(v)
}
return res
}



/*:52*/



/*55:*/


//line amail.w:762

func(this messages)Search(id int)(int,bool){
pos:=sort.Search(len(this),func(i int)bool{return this[i].id<=id});
if pos!=len(this)&&this[pos].id==id{
return pos,true
}
return pos,false
}



/*:55*/



/*56:*/


//line amail.w:772

func(this*messages)Insert(msg*message,pos int){
*this= append(*this,nil)
copy((*this)[pos+1:],(*this)[pos:])
(*this)[pos]= msg
}



/*:56*/



/*57:*/


//line amail.w:782

func(this*messages)SearchInsert(msg*message)(int,bool){
pos,ok:=this.Search(msg.id)
if ok{
return pos,false
}
this.Insert(msg,pos)
return pos,true
}



/*:57*/



/*58:*/


//line amail.w:795

func(this*messages)Delete(pos int)(*message,bool){
if pos<0||pos> len(*this)-1{
return nil,false
}
msg:=(*this)[pos]
*this= append((*this)[:pos],(*this)[pos+1:]...)
return msg,true
}



/*:58*/



/*59:*/


//line amail.w:808

func(this*messages)DeleteById(id int)(*message,bool){
pos,ok:=this.Search(id)
if!ok{
return nil,false
}
return this.Delete(pos)
}



/*:59*/



/*61:*/


//line amail.w:822

func(box*mailbox)loop(){
glog.V(debug).Infof("start a message loop for the '%s' mailbox\n",box.name)
counted:=false
pcount:=0
ontop:=false


/*49:*/


//line amail.w:666

{
glog.V(debug).Infof("counting of messages in the '%s' mailbox\n",box.name)
var err error
box.fid,err= rfid.Walk(box.name)
if err!=nil{
glog.Errorf("can't walk to '%s': %v",box.name,err)
return
}
defer box.fid.Close()
fs,err:=box.fid.Dirreadall()
if err!=nil{
glog.Errorf("can't read a mailbox '%s': %s",box.name,err)
return
}
box.total= len(fs)-2
box.all= make(messages,0,box.total)
for i:=len(fs)-1;i>=0;{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:685



/*62:*/


//line amail.w:846

case id:=<-box.mch:
glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n",id,box.name)
msg,new,err:=box.newMessage(id)
if err!=nil{
continue
}
if new{


/*65:*/


//line amail.w:884

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:65*/


//line amail.w:854

}
box.total++


/*66:*/


//line amail.w:892

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:66*/


//line amail.w:857



/*149:*/


//line amail.w:1847

{
glog.V(debug).Infof("print '%s/%d' at exact position\n",box.name,msg.id)
mrfch<-&refresh{seek|insert|exact,append(messages{},msg)}
}



/*:149*/


//line amail.w:858

if!box.thread{
if box.threadMode(){


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:861

var msgs messages
src:=append(messages{},msg)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:864



/*145:*/


//line amail.w:1814

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)
box.rfch<-&refresh{seek|insert,msgs}
}



/*:145*/


//line amail.w:865

}else{


/*144:*/


//line amail.w:1807

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&refresh{seek|insert,append(messages{},msg)}
}



/*:144*/


//line amail.w:867

}
}


/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:870





/*:62*/



/*63:*/


//line amail.w:874

case id:=<-box.dch:
glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n",id,box.name)


/*67:*/


//line amail.w:900

if i,ok:=box.all.Search(id);ok{
msgs:=append(messages{},box.all[i])


/*68:*/


//line amail.w:909

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*132:*/


//line amail.w:1673

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{id string;val interface{}}{id:msg.messageid}
}



/*:132*/


//line amail.w:918



/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:919

}
}



/*:68*/


//line amail.w:903



/*112:*/


//line amail.w:1444

mdch<-msgs



/*:112*/


//line amail.w:904

}



/*:67*/


//line amail.w:877




/*:63*/



/*75:*/


//line amail.w:967

case<-box.cch:
glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
box.shownew= shownew
box.showthreads= showthreads
box.thread= false


/*76:*/


//line amail.w:983

glog.V(debug).Infof("creation a window for the '%s' mailbox\n",box.name)
var err error
if box.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
if box.ech,err= box.w.EventChannel(0,goacme.Mouse,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*167:*/


//line amail.w:2044

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2047




/*:167*/


//line amail.w:994



/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:995





/*:76*/


//line amail.w:974



/*140:*/


//line amail.w:1753

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1756

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1763

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1765

}
}
box.rfch<-&refresh{0,msgs}
}



/*:140*/


//line amail.w:975



/*8:*/


//line amail.w:110

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:976

}else{
glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n",box.name)
box.w.WriteCtl("dot=addr\nshow")
}



/*:75*/



/*78:*/


//line amail.w:1003

case ev,ok:=<-box.ech:
glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n",box.name,ev)
if!ok{
box.ech= nil
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":


/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:1013

box.w.Del(true)
box.w.Close()
box.w= nil


/*9:*/


//line amail.w:115

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:1017

continue
case"ShowNew":
box.thread= false
box.shownew= true
case"ShowAll":
box.thread= false
box.shownew= false
case"ShowThreads":
box.showthreads= true
if box.shownew==true{


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1028

continue
}
case"ShowPlain":
box.showthreads= false
if box.shownew==true{


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1034

continue
}
case"Thread":
var msg*message
if len(ev.Arg)==0{


/*88:*/


//line amail.w:1181

glog.V(debug).Infof("getting a pointer to current message in the window of the '%s' mailbox\n",box.name)
num:=0
if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else if err:=box.w.WriteAddr("-/^/");err!=nil{
glog.V(debug).Infof("can't write to 'addr': %v\n",err)
}else if err:=box.w.WriteAddr("/[0-9]+(%s)?\\//",escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if str,err:=bufio.NewReader(data).ReadString('/');err!=nil{
glog.Errorf("can't read from 'data' file: %s\n",err)
}else if _,err:=fmt.Sscanf(strings.TrimLeft(str,levelmark),"%d",&num);err==nil{
glog.V(debug).Infof("current message: %d(%s)\n",num,str)
if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
}else{
glog.V(debug).Infof("can't get a current message from: %s\n",str)
}



/*:88*/


//line amail.w:1040

}else if num,err:=strconv.Atoi(strings.TrimSpace(ev.Arg));err!=nil{
continue
}else if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
if msg!=nil{
box.thread= true


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1048



/*169:*/


//line amail.w:2084

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:169*/


//line amail.w:1049



/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:1050



/*143:*/


//line amail.w:1799



/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1800

var msgs messages
src:=append(messages{},msg)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1803

box.rfch<-&refresh{0,msgs}



/*:143*/


//line amail.w:1051

}
continue
case"Delmesg":


/*95:*/


//line amail.w:1255



/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:1256

if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*80:*/


//line amail.w:1111

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1117

if err==io.EOF{
break
}
}
}




/*:80*/


//line amail.w:1260

if len(msgs)!=0{


/*99:*/


//line amail.w:1305

glog.V(debug).Infoln("sending messages to mark for deletion")
markch<-&msgs



/*:99*/


//line amail.w:1262

continue
}



/*:95*/


//line amail.w:1055

continue
case"UnDelmesg":


/*96:*/


//line amail.w:1267



/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:1268

if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*80:*/


//line amail.w:1111

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1117

if err==io.EOF{
break
}
}
}




/*:80*/


//line amail.w:1272

if len(msgs)!=0{


/*100:*/


//line amail.w:1310

glog.V(debug).Infoln("sending messages to unmark for deletion")
unmarkch<-&msgs




/*:100*/


//line amail.w:1274

continue
}




/*:96*/


//line amail.w:1058

continue
case"Put":


/*106:*/


//line amail.w:1375

f,err:=box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl': %v\n",err)
continue
}
var msgs messages
for i:=0;i<len(box.all);{
if!box.all[i].deleted||box.all[i].w!=nil{
i++
continue
}
msgs= append(msgs,box.all[i])


/*68:*/


//line amail.w:909

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*132:*/


//line amail.w:1673

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{id string;val interface{}}{id:msg.messageid}
}



/*:132*/


//line amail.w:918



/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:919

}
}



/*:68*/


//line amail.w:1391

}
cmd:=fmt.Sprintf("delete %s",box.name)
for _,msg:=range msgs{
cmd= fmt.Sprintf("%s %d ",cmd,msg.id)
}
glog.V(debug).Infof("command to delete messages: '%s'\n",cmd)
if _,err:=f.Write([]byte(cmd));err!=nil{
glog.Errorf("can't delete messages: %v\n",err)
}
f.Close()


/*112:*/


//line amail.w:1444

mdch<-msgs



/*:112*/


//line amail.w:1402





/*:106*/


//line amail.w:1061

continue
case"Mail":
var msg*message


/*224:*/


//line amail.w:3063

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*225:*/


//line amail.w:3075

go func(msg*message){
glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
for ev,err:=w.ReadEvent();err==nil;ev,err= w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
w.UnreadEvent(ev)
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
w.UnreadEvent(ev)
w.Close()
return
case"Post":


/*229:*/


//line amail.w:3162

{


/*231:*/


//line amail.w:3279



/*221:*/


//line amail.w:2994

once.Do(func(){

/*209:*/


//line amail.w:2810

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:209*/



/*218:*/


//line amail.w:2977

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:218*/



/*232:*/


//line amail.w:3283

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:232*/


//line amail.w:2995
})



/*:221*/


//line amail.w:3280




/*:231*/


//line amail.w:3164

w.Seek(0,0)
w.WriteAddr(wholefile)
ff,_:=w.File("xdata")
b:=bufio.NewReader(ff)
var to,cc,bcc,attach,include[]string
var subject string
for{
s,err:=b.ReadString('\n')
if err!=nil{
break
}
s= strings.TrimSpace(s)
if len(s)==0{
// an empty line, the rest is a body of the message
break
}
p:=strings.Index(s,":")
if p!=-1{
f:=strings.Split(s[p+1:],",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
switch strings.ToLower(s[:p]){
case"from","to":
to= append(to,f...)
case"cc":
cc= append(cc,f...)
case"bcc":
bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%q",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
to= append(to,f...)
}
}
args:=append([]string{},"-8")
if msg!=nil{
args= append(args,"-R",fmt.Sprintf("%s/%d",msg.box.name,msg.id))
}
if len(subject)!=0{
args= append(args,"-s",subject)
}
for _,v:=range include{
args= append(args,"-A",v)
}
for _,v:=range attach{
args= append(args,"-a",v)
}

c:=exec.Command(plan9dir+"/bin/upas/marshal",args...)
p,err:=c.StdinPipe()
if err!=nil{
glog.Errorf("can't get a stdin pipe: %v\n",err)
continue
}
if err:=c.Start();err!=nil{
glog.Errorf("can't start 'upas/marshal': %v\n",err)
continue
}
if len(to)!=0{
if _,err:=fmt.Fprintln(p,"To:",strings.Join(to,","));err!=nil{
glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("to is written")
if len(cc)!=0{
if _,err:=fmt.Fprintln(p,"CC:",strings.Join(cc,","));err!=nil{
glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("cc is written")
if len(bcc)!=0{
if _,err:=fmt.Fprintln(p,"BCC:",strings.Join(bcc,","));err!=nil{
glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("bcc is written")
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("writing '%s':%v",s,err)

p.Write([]byte(s))
if err==io.EOF{
p.Write([]byte("\n"))
break
}
}
glog.V(debug).Infoln("body is written")
p.Write([]byte("\n"))
p.Close()
c.Wait()
w.Del(true)
w.Close()
}





/*:229*/


//line amail.w:3090

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:225*/


//line amail.w:3072




/*:224*/


//line amail.w:1065

name:=fmt.Sprintf("Amail/%s/New",box.name)


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1067

continue
case"Search":
glog.V(debug).Infof("search argument: '%s'\n",ev.Arg)


/*176:*/


//line amail.w:2195

{
msgs:=box.search(ev.Arg)


/*169:*/


//line amail.w:2084

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:169*/


//line amail.w:2198



/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:2199

name:=fmt.Sprintf("Amail/%s/Search(%s)",box.name,strings.Replace(ev.Arg," ","␣",-1))
w:=box.w
box.thread= false
box.shownew= true
box.showthreads= false


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2205

glog.V(debug).Infof("len of msgs: %v\n",len(msgs))
box.rfch<-&refresh{0,msgs}
}



/*:176*/


//line amail.w:1071

continue
default:
box.w.UnreadEvent(ev)
continue
}


/*167:*/


//line amail.w:2044

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2047




/*:167*/


//line amail.w:1077



/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1078



/*169:*/


//line amail.w:2084

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:169*/


//line amail.w:1079



/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:1080



/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:1081



/*140:*/


//line amail.w:1753

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1756

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1763

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1765

}
}
box.rfch<-&refresh{0,msgs}
}



/*:140*/


//line amail.w:1082

continue
}else if(ev.Type&goacme.Look)==goacme.Look{


/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:1085

if(ev.Type&goacme.Tag)==goacme.Tag{
s:=ev.Text


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1088

}else{


/*79:*/


//line amail.w:1103

glog.V(debug).Infof("event: %v\n",ev)
if err:=box.w.WriteAddr("#%d,#%d",ev.Begin,ev.End);err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else


/*80:*/


//line amail.w:1111

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1117

if err==io.EOF{
break
}
}
}




/*:80*/


//line amail.w:1108




/*:79*/


//line amail.w:1090

}
if len(msgs)!=0{


/*87:*/


//line amail.w:1177

lch<-&msgs



/*:87*/


//line amail.w:1093

continue
}
}
box.w.UnreadEvent(ev)



/*:78*/



/*103:*/


//line amail.w:1326

case ids:=<-box.markch:
var msgs messages
for _,id:=range ids{


/*104:*/


//line amail.w:1341

if p,ok:=box.all.Search(id);ok{
if box.all[p].deleted{
continue
}
box.all[p].deleted= true
box.deleted++
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*191:*/


//line amail.w:2408

msg.writeTag()



/*:191*/


//line amail.w:1351

}
glog.V(debug).Infof("the '%v' message is marked for deletion\n",id)
}



/*:104*/


//line amail.w:1330

}


/*148:*/


//line amail.w:1838

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:148*/


//line amail.w:1332

case ids:=<-box.unmarkch:
var msgs messages
for _,id:=range ids{


/*105:*/


//line amail.w:1357

if p,ok:=box.all.Search(id);ok{
if!box.all[p].deleted{
continue
}
box.all[p].deleted= false
box.deleted--
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*191:*/


//line amail.w:2408

msg.writeTag()



/*:191*/


//line amail.w:1367

}
glog.V(debug).Infof("the '%v' message is unmarked for deletion\n",id)

}



/*:105*/


//line amail.w:1336

}


/*148:*/


//line amail.w:1838

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:148*/


//line amail.w:1338




/*:103*/



/*109:*/


//line amail.w:1415

case m:=<-box.mdch:
if box.w==nil{
continue
}
glog.V(debug).Infof("%d messages were received to be deleted from the '%s' mailbox\n",len(m),box.name)
for _,msg:=range m{


/*185:*/


//line amail.w:2318

box.eraseMessage(msg)




/*:185*/


//line amail.w:1422

if box.threadMode(){


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:1424



/*163:*/


//line amail.w:1996

{
if len(children)!=0{
var msgs messages
var src messages
for _,msg:=range children{
if msg.box!=box{


/*185:*/


//line amail.w:2318

box.eraseMessage(msg)




/*:185*/


//line amail.w:2003

}else{


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:2005

}
}


/*147:*/


//line amail.w:1829

{
if len(msgs)!=0{
glog.V(debug).Infof("inform the '%s' mailbox to refresh messages\n",box.name)
box.rfch<-&refresh{seek,msgs}
}
}



/*:147*/


//line amail.w:2008

}
}



/*:163*/


//line amail.w:1425

}
}


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1428




/*:109*/



/*138:*/


//line amail.w:1725

case v:=<-box.rfch:
box.irfch<-v

case v:=<-box.irfch:
glog.V(debug).Infof("a signal to print message of the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
glog.V(debug).Infof("a window of the '%s' mailbox doesn't exist, ignore the signal\n",box.name)
continue
}
if box.threadMode()&&!counted{
glog.V(debug).Infof("counting of threads of the '%s' mailbox is not finished, ignore the signal\n",box.name)
continue
}


/*157:*/


//line amail.w:1907

{
glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with flags: %v\n",box.name,len(v.msgs),v.flags)
f,err:=box.w.File("data")
if err!=nil{
glog.Errorf("can't open 'data' file of the '%s' messgebox: %v\n",box.name,err)
continue
}
if(v.flags&seek)==seek{


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1916

msg:=v.msgs[0]


/*172:*/


//line amail.w:2107



/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:2108



/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2109

glog.V(debug).Infof("composed message addr '%s' in the '%s' mailbox\n",addr,box.name)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("the '%d' message is not found in the window of the '%s' mailbox\n",msg.id,box.name)
if(v.flags&insert)==0{


/*173:*/


//line amail.w:2139

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:2142

continue




/*:173*/


//line amail.w:2114

}
if box.threadMode(){


/*175:*/


//line amail.w:2158

if msg.parent!=nil{
glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n",msg.id)
m:=msg
msg= m.parent
found:=false
for!found{


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:2165

if len(children)==0{
break
}
for i,v:=range children{
if v==m{
if i==0{
found= true
}
break
}
msg= v
}
}
glog.V(debug).Infof("the '%d' message will be printed after the '%d' message\n",m.id,msg.id)


/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2180

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
if(v.flags&exact)==exact{


/*173:*/


//line amail.w:2139

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:2142

continue




/*:173*/


//line amail.w:2185

}
}
}else if(v.flags&exact)==exact{


/*173:*/


//line amail.w:2139

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:2142

continue




/*:173*/


//line amail.w:2189

}else if err:=box.w.WriteAddr(bof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %v\n",bof,box.name,err)
}



/*:175*/


//line amail.w:2117

}else if p,ok:=src.Search(msg.id);!ok{
glog.V(debug).Infof("the '%d' message is not found in  the '%s' mailbox's window\n",msg.id,box.name)
}else if p==0{
if err:=box.w.WriteAddr(bof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %s\n",bof,box.name,err)
}
}else if p==len(src)-1{
if err:=box.w.WriteAddr(eof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %s\n",eof,box.name,err)
}
}else{
msg:=src[p-1]


/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2130

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
}
}
}



/*:172*/


//line amail.w:1918

}else if err:=box.w.WriteAddr(eof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",eof,err)
continue
}
w:=box.w
glog.V(debug).Infof("printing of messages of the '%s' mailbox\n",box.name)
buf:=make([]byte,0,0x8000)


/*160:*/


//line amail.w:1958

c:=0
for len(v.msgs)> 0&&c<100{
msg:=v.msgs[0]
glog.V(debug).Infof("printing of '%s/%d' message with in the '%s' mailbox\n",msg.box.name,msg.id,box.name)
if box.threadMode(){


/*162:*/


//line amail.w:1986

{
for p:=msg.parent;p!=nil;p= p.parent{
buf= append(buf,levelmark...)
}
}



/*:162*/


//line amail.w:1964

}
c++


/*90:*/


//line amail.w:1208

glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n",msg.id,box.name)
buf= append(buf,fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n",
func()string{if msg.deleted{return deleted};return""}(),
func()string{if msg.box!=box{return fmt.Sprintf("%s/",msg.box.name)};return""}(),
msg.id,
func()string{if msg.unread{return newmark};return""}(),
msg.from,
msg.date,
msg.subject)...)




/*:90*/


//line amail.w:1967

v.msgs= v.msgs[1:]
if(v.flags&seek)==seek{
break
}
}
pcount+= c



/*:160*/


//line amail.w:1926

if _,err:=f.Write(buf);err!=nil{
glog.Errorf("can't write to 'data' file of the '%s' messgebox: %v\n",box.name,err)
}


/*159:*/


//line amail.w:1944

if!ontop{
glog.V(debug).Infof("pcount: %v, ontop: %v\n",pcount,ontop)


/*94:*/


//line amail.w:1249

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:94*/


//line amail.w:1947

if pcount>=100{
ontop= true
}
}




/*:159*/


//line amail.w:1930



/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:1931

}



/*:157*/


//line amail.w:1739





/*:138*/



/*181:*/


//line amail.w:2261

case ids:=<-box.lch:
var msgs messages
for _,id:=range ids{
glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n",id,box.name)
p,ok:=box.all.Search(id)
if!ok{
glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n",id,box.name)
continue
}
msg:=box.all[p]
if msg.w==nil{
if msg.unread{


/*182:*/


//line amail.w:2289

msg.unread= false
box.unread.DeleteById(id)



/*:182*/


//line amail.w:2274



/*183:*/


//line amail.w:2295

if!box.thread&&box.shownew{


/*185:*/


//line amail.w:2318

box.eraseMessage(msg)




/*:185*/


//line amail.w:2297



/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:2298

}
msgs= append(msgs,msg)





/*:183*/


//line amail.w:2275



/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:2276

}
if err:=msg.open();err!=nil{
continue
}
}else{
glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n",id,box.name)
msg.w.WriteCtl("dot=addr\nshow")
}
}


/*148:*/


//line amail.w:1838

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:148*/


//line amail.w:2286




/*:181*/


//line amail.w:686

default:
d:=fs[i]
i--
if(d.Mode&plan9.DMDIR)!=plan9.DMDIR{
continue
}
id,err:=strconv.Atoi(d.Name)
if err!=nil{// it seems this is a mailbox
// decrease a total number of messages
box.total--
name:=box.name+"/"+d.Name


/*27:*/


//line amail.w:304

glog.V(debug).Infof("send a mailbox '%s' to put in the list\n",name)
bch<-name



/*:27*/


//line amail.w:698

continue
}
if msg,new,_:=box.newMessage(id);err==nil{
if new{


/*65:*/


//line amail.w:884

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:65*/


//line amail.w:703

}


/*66:*/


//line amail.w:892

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:66*/


//line amail.w:705

}else{
glog.Errorf("can't create a new '%d' message in the '%s' mailbox\n",id,box.name)
}


/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:709



/*154:*/


//line amail.w:1879

if!box.threadMode(){


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1881

if len(src)!=0&&len(src)%500==0{
glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n",box.name,len(src)-box.pos)
msgs:=append(messages{},src[box.pos:len(src)]...)
box.pos= len(src)
box.rfch<-&refresh{0,msgs}
}
}



/*:154*/


//line amail.w:710

}
}


/*155:*/


//line amail.w:1891

if!box.threadMode(){


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1893

if box.pos!=len(src){
glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n",box.name,len(src)-box.pos)
msgs:=append(messages{},src[box.pos:len(src)]...)
box.pos= len(src)
box.rfch<-&refresh{0,msgs}
}
}



/*:155*/


//line amail.w:713

}



/*:49*/


//line amail.w:828

counted= true
if box.threadMode(){


/*140:*/


//line amail.w:1753

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1756

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1763

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1765

}
}
box.rfch<-&refresh{0,msgs}
}



/*:140*/


//line amail.w:831

}
defer glog.V(debug).Infof("a message loop of the '%s' mailbox is done\n",box.name)
for{
select{


/*4:*/


//line amail.w:82

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:836



/*62:*/


//line amail.w:846

case id:=<-box.mch:
glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n",id,box.name)
msg,new,err:=box.newMessage(id)
if err!=nil{
continue
}
if new{


/*65:*/


//line amail.w:884

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:65*/


//line amail.w:854

}
box.total++


/*66:*/


//line amail.w:892

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:66*/


//line amail.w:857



/*149:*/


//line amail.w:1847

{
glog.V(debug).Infof("print '%s/%d' at exact position\n",box.name,msg.id)
mrfch<-&refresh{seek|insert|exact,append(messages{},msg)}
}



/*:149*/


//line amail.w:858

if!box.thread{
if box.threadMode(){


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:861

var msgs messages
src:=append(messages{},msg)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:864



/*145:*/


//line amail.w:1814

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)
box.rfch<-&refresh{seek|insert,msgs}
}



/*:145*/


//line amail.w:865

}else{


/*144:*/


//line amail.w:1807

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&refresh{seek|insert,append(messages{},msg)}
}



/*:144*/


//line amail.w:867

}
}


/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:870





/*:62*/



/*63:*/


//line amail.w:874

case id:=<-box.dch:
glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n",id,box.name)


/*67:*/


//line amail.w:900

if i,ok:=box.all.Search(id);ok{
msgs:=append(messages{},box.all[i])


/*68:*/


//line amail.w:909

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*132:*/


//line amail.w:1673

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{id string;val interface{}}{id:msg.messageid}
}



/*:132*/


//line amail.w:918



/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:919

}
}



/*:68*/


//line amail.w:903



/*112:*/


//line amail.w:1444

mdch<-msgs



/*:112*/


//line amail.w:904

}



/*:67*/


//line amail.w:877




/*:63*/



/*75:*/


//line amail.w:967

case<-box.cch:
glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
box.shownew= shownew
box.showthreads= showthreads
box.thread= false


/*76:*/


//line amail.w:983

glog.V(debug).Infof("creation a window for the '%s' mailbox\n",box.name)
var err error
if box.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
if box.ech,err= box.w.EventChannel(0,goacme.Mouse,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*167:*/


//line amail.w:2044

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2047




/*:167*/


//line amail.w:994



/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:995





/*:76*/


//line amail.w:974



/*140:*/


//line amail.w:1753

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1756

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1763

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1765

}
}
box.rfch<-&refresh{0,msgs}
}



/*:140*/


//line amail.w:975



/*8:*/


//line amail.w:110

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:976

}else{
glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n",box.name)
box.w.WriteCtl("dot=addr\nshow")
}



/*:75*/



/*78:*/


//line amail.w:1003

case ev,ok:=<-box.ech:
glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n",box.name,ev)
if!ok{
box.ech= nil
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":


/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:1013

box.w.Del(true)
box.w.Close()
box.w= nil


/*9:*/


//line amail.w:115

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:1017

continue
case"ShowNew":
box.thread= false
box.shownew= true
case"ShowAll":
box.thread= false
box.shownew= false
case"ShowThreads":
box.showthreads= true
if box.shownew==true{


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1028

continue
}
case"ShowPlain":
box.showthreads= false
if box.shownew==true{


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1034

continue
}
case"Thread":
var msg*message
if len(ev.Arg)==0{


/*88:*/


//line amail.w:1181

glog.V(debug).Infof("getting a pointer to current message in the window of the '%s' mailbox\n",box.name)
num:=0
if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else if err:=box.w.WriteAddr("-/^/");err!=nil{
glog.V(debug).Infof("can't write to 'addr': %v\n",err)
}else if err:=box.w.WriteAddr("/[0-9]+(%s)?\\//",escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if str,err:=bufio.NewReader(data).ReadString('/');err!=nil{
glog.Errorf("can't read from 'data' file: %s\n",err)
}else if _,err:=fmt.Sscanf(strings.TrimLeft(str,levelmark),"%d",&num);err==nil{
glog.V(debug).Infof("current message: %d(%s)\n",num,str)
if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
}else{
glog.V(debug).Infof("can't get a current message from: %s\n",str)
}



/*:88*/


//line amail.w:1040

}else if num,err:=strconv.Atoi(strings.TrimSpace(ev.Arg));err!=nil{
continue
}else if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
if msg!=nil{
box.thread= true


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1048



/*169:*/


//line amail.w:2084

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:169*/


//line amail.w:1049



/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:1050



/*143:*/


//line amail.w:1799



/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1800

var msgs messages
src:=append(messages{},msg)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1803

box.rfch<-&refresh{0,msgs}



/*:143*/


//line amail.w:1051

}
continue
case"Delmesg":


/*95:*/


//line amail.w:1255



/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:1256

if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*80:*/


//line amail.w:1111

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1117

if err==io.EOF{
break
}
}
}




/*:80*/


//line amail.w:1260

if len(msgs)!=0{


/*99:*/


//line amail.w:1305

glog.V(debug).Infoln("sending messages to mark for deletion")
markch<-&msgs



/*:99*/


//line amail.w:1262

continue
}



/*:95*/


//line amail.w:1055

continue
case"UnDelmesg":


/*96:*/


//line amail.w:1267



/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:1268

if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*80:*/


//line amail.w:1111

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1117

if err==io.EOF{
break
}
}
}




/*:80*/


//line amail.w:1272

if len(msgs)!=0{


/*100:*/


//line amail.w:1310

glog.V(debug).Infoln("sending messages to unmark for deletion")
unmarkch<-&msgs




/*:100*/


//line amail.w:1274

continue
}




/*:96*/


//line amail.w:1058

continue
case"Put":


/*106:*/


//line amail.w:1375

f,err:=box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl': %v\n",err)
continue
}
var msgs messages
for i:=0;i<len(box.all);{
if!box.all[i].deleted||box.all[i].w!=nil{
i++
continue
}
msgs= append(msgs,box.all[i])


/*68:*/


//line amail.w:909

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*132:*/


//line amail.w:1673

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{id string;val interface{}}{id:msg.messageid}
}



/*:132*/


//line amail.w:918



/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:919

}
}



/*:68*/


//line amail.w:1391

}
cmd:=fmt.Sprintf("delete %s",box.name)
for _,msg:=range msgs{
cmd= fmt.Sprintf("%s %d ",cmd,msg.id)
}
glog.V(debug).Infof("command to delete messages: '%s'\n",cmd)
if _,err:=f.Write([]byte(cmd));err!=nil{
glog.Errorf("can't delete messages: %v\n",err)
}
f.Close()


/*112:*/


//line amail.w:1444

mdch<-msgs



/*:112*/


//line amail.w:1402





/*:106*/


//line amail.w:1061

continue
case"Mail":
var msg*message


/*224:*/


//line amail.w:3063

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*225:*/


//line amail.w:3075

go func(msg*message){
glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
for ev,err:=w.ReadEvent();err==nil;ev,err= w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
w.UnreadEvent(ev)
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
w.UnreadEvent(ev)
w.Close()
return
case"Post":


/*229:*/


//line amail.w:3162

{


/*231:*/


//line amail.w:3279



/*221:*/


//line amail.w:2994

once.Do(func(){

/*209:*/


//line amail.w:2810

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:209*/



/*218:*/


//line amail.w:2977

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:218*/



/*232:*/


//line amail.w:3283

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:232*/


//line amail.w:2995
})



/*:221*/


//line amail.w:3280




/*:231*/


//line amail.w:3164

w.Seek(0,0)
w.WriteAddr(wholefile)
ff,_:=w.File("xdata")
b:=bufio.NewReader(ff)
var to,cc,bcc,attach,include[]string
var subject string
for{
s,err:=b.ReadString('\n')
if err!=nil{
break
}
s= strings.TrimSpace(s)
if len(s)==0{
// an empty line, the rest is a body of the message
break
}
p:=strings.Index(s,":")
if p!=-1{
f:=strings.Split(s[p+1:],",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
switch strings.ToLower(s[:p]){
case"from","to":
to= append(to,f...)
case"cc":
cc= append(cc,f...)
case"bcc":
bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%q",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
to= append(to,f...)
}
}
args:=append([]string{},"-8")
if msg!=nil{
args= append(args,"-R",fmt.Sprintf("%s/%d",msg.box.name,msg.id))
}
if len(subject)!=0{
args= append(args,"-s",subject)
}
for _,v:=range include{
args= append(args,"-A",v)
}
for _,v:=range attach{
args= append(args,"-a",v)
}

c:=exec.Command(plan9dir+"/bin/upas/marshal",args...)
p,err:=c.StdinPipe()
if err!=nil{
glog.Errorf("can't get a stdin pipe: %v\n",err)
continue
}
if err:=c.Start();err!=nil{
glog.Errorf("can't start 'upas/marshal': %v\n",err)
continue
}
if len(to)!=0{
if _,err:=fmt.Fprintln(p,"To:",strings.Join(to,","));err!=nil{
glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("to is written")
if len(cc)!=0{
if _,err:=fmt.Fprintln(p,"CC:",strings.Join(cc,","));err!=nil{
glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("cc is written")
if len(bcc)!=0{
if _,err:=fmt.Fprintln(p,"BCC:",strings.Join(bcc,","));err!=nil{
glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("bcc is written")
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("writing '%s':%v",s,err)

p.Write([]byte(s))
if err==io.EOF{
p.Write([]byte("\n"))
break
}
}
glog.V(debug).Infoln("body is written")
p.Write([]byte("\n"))
p.Close()
c.Wait()
w.Del(true)
w.Close()
}





/*:229*/


//line amail.w:3090

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:225*/


//line amail.w:3072




/*:224*/


//line amail.w:1065

name:=fmt.Sprintf("Amail/%s/New",box.name)


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1067

continue
case"Search":
glog.V(debug).Infof("search argument: '%s'\n",ev.Arg)


/*176:*/


//line amail.w:2195

{
msgs:=box.search(ev.Arg)


/*169:*/


//line amail.w:2084

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:169*/


//line amail.w:2198



/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:2199

name:=fmt.Sprintf("Amail/%s/Search(%s)",box.name,strings.Replace(ev.Arg," ","␣",-1))
w:=box.w
box.thread= false
box.shownew= true
box.showthreads= false


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2205

glog.V(debug).Infof("len of msgs: %v\n",len(msgs))
box.rfch<-&refresh{0,msgs}
}



/*:176*/


//line amail.w:1071

continue
default:
box.w.UnreadEvent(ev)
continue
}


/*167:*/


//line amail.w:2044

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2047




/*:167*/


//line amail.w:1077



/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1078



/*169:*/


//line amail.w:2084

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:169*/


//line amail.w:1079



/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:1080



/*153:*/


//line amail.w:1873

box.pos= 0
ontop= false



/*:153*/



/*161:*/


//line amail.w:1976

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:161*/


//line amail.w:1081



/*140:*/


//line amail.w:1753

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:1756

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*120:*/


//line amail.w:1507

for msg.parent!=nil{
msg= msg.parent
}



/*:120*/


//line amail.w:1763

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:1765

}
}
box.rfch<-&refresh{0,msgs}
}



/*:140*/


//line amail.w:1082

continue
}else if(ev.Type&goacme.Look)==goacme.Look{


/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:1085

if(ev.Type&goacme.Tag)==goacme.Tag{
s:=ev.Text


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1088

}else{


/*79:*/


//line amail.w:1103

glog.V(debug).Infof("event: %v\n",ev)
if err:=box.w.WriteAddr("#%d,#%d",ev.Begin,ev.End);err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else


/*80:*/


//line amail.w:1111

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*81:*/


//line amail.w:1127

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
id:=0
for i,v:=range f{
var err error
if id,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",id,name)
}


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:1142

break
}
}
}



/*:81*/


//line amail.w:1117

if err==io.EOF{
break
}
}
}




/*:80*/


//line amail.w:1108




/*:79*/


//line amail.w:1090

}
if len(msgs)!=0{


/*87:*/


//line amail.w:1177

lch<-&msgs



/*:87*/


//line amail.w:1093

continue
}
}
box.w.UnreadEvent(ev)



/*:78*/



/*103:*/


//line amail.w:1326

case ids:=<-box.markch:
var msgs messages
for _,id:=range ids{


/*104:*/


//line amail.w:1341

if p,ok:=box.all.Search(id);ok{
if box.all[p].deleted{
continue
}
box.all[p].deleted= true
box.deleted++
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*191:*/


//line amail.w:2408

msg.writeTag()



/*:191*/


//line amail.w:1351

}
glog.V(debug).Infof("the '%v' message is marked for deletion\n",id)
}



/*:104*/


//line amail.w:1330

}


/*148:*/


//line amail.w:1838

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:148*/


//line amail.w:1332

case ids:=<-box.unmarkch:
var msgs messages
for _,id:=range ids{


/*105:*/


//line amail.w:1357

if p,ok:=box.all.Search(id);ok{
if!box.all[p].deleted{
continue
}
box.all[p].deleted= false
box.deleted--
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*191:*/


//line amail.w:2408

msg.writeTag()



/*:191*/


//line amail.w:1367

}
glog.V(debug).Infof("the '%v' message is unmarked for deletion\n",id)

}



/*:105*/


//line amail.w:1336

}


/*148:*/


//line amail.w:1838

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:148*/


//line amail.w:1338




/*:103*/



/*109:*/


//line amail.w:1415

case m:=<-box.mdch:
if box.w==nil{
continue
}
glog.V(debug).Infof("%d messages were received to be deleted from the '%s' mailbox\n",len(m),box.name)
for _,msg:=range m{


/*185:*/


//line amail.w:2318

box.eraseMessage(msg)




/*:185*/


//line amail.w:1422

if box.threadMode(){


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:1424



/*163:*/


//line amail.w:1996

{
if len(children)!=0{
var msgs messages
var src messages
for _,msg:=range children{
if msg.box!=box{


/*185:*/


//line amail.w:2318

box.eraseMessage(msg)




/*:185*/


//line amail.w:2003

}else{


/*141:*/


//line amail.w:1773

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:141*/


//line amail.w:2005

}
}


/*147:*/


//line amail.w:1829

{
if len(msgs)!=0{
glog.V(debug).Infof("inform the '%s' mailbox to refresh messages\n",box.name)
box.rfch<-&refresh{seek,msgs}
}
}



/*:147*/


//line amail.w:2008

}
}



/*:163*/


//line amail.w:1425

}
}


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1428




/*:109*/



/*138:*/


//line amail.w:1725

case v:=<-box.rfch:
box.irfch<-v

case v:=<-box.irfch:
glog.V(debug).Infof("a signal to print message of the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
glog.V(debug).Infof("a window of the '%s' mailbox doesn't exist, ignore the signal\n",box.name)
continue
}
if box.threadMode()&&!counted{
glog.V(debug).Infof("counting of threads of the '%s' mailbox is not finished, ignore the signal\n",box.name)
continue
}


/*157:*/


//line amail.w:1907

{
glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with flags: %v\n",box.name,len(v.msgs),v.flags)
f,err:=box.w.File("data")
if err!=nil{
glog.Errorf("can't open 'data' file of the '%s' messgebox: %v\n",box.name,err)
continue
}
if(v.flags&seek)==seek{


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:1916

msg:=v.msgs[0]


/*172:*/


//line amail.w:2107



/*139:*/


//line amail.w:1743

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:139*/


//line amail.w:2108



/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2109

glog.V(debug).Infof("composed message addr '%s' in the '%s' mailbox\n",addr,box.name)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("the '%d' message is not found in the window of the '%s' mailbox\n",msg.id,box.name)
if(v.flags&insert)==0{


/*173:*/


//line amail.w:2139

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:2142

continue




/*:173*/


//line amail.w:2114

}
if box.threadMode(){


/*175:*/


//line amail.w:2158

if msg.parent!=nil{
glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n",msg.id)
m:=msg
msg= m.parent
found:=false
for!found{


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:2165

if len(children)==0{
break
}
for i,v:=range children{
if v==m{
if i==0{
found= true
}
break
}
msg= v
}
}
glog.V(debug).Infof("the '%d' message will be printed after the '%d' message\n",m.id,msg.id)


/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2180

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
if(v.flags&exact)==exact{


/*173:*/


//line amail.w:2139

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:2142

continue




/*:173*/


//line amail.w:2185

}
}
}else if(v.flags&exact)==exact{


/*173:*/


//line amail.w:2139

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:2142

continue




/*:173*/


//line amail.w:2189

}else if err:=box.w.WriteAddr(bof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %v\n",bof,box.name,err)
}



/*:175*/


//line amail.w:2117

}else if p,ok:=src.Search(msg.id);!ok{
glog.V(debug).Infof("the '%d' message is not found in  the '%s' mailbox's window\n",msg.id,box.name)
}else if p==0{
if err:=box.w.WriteAddr(bof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %s\n",bof,box.name,err)
}
}else if p==len(src)-1{
if err:=box.w.WriteAddr(eof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %s\n",eof,box.name,err)
}
}else{
msg:=src[p-1]


/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2130

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
}
}
}



/*:172*/


//line amail.w:1918

}else if err:=box.w.WriteAddr(eof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",eof,err)
continue
}
w:=box.w
glog.V(debug).Infof("printing of messages of the '%s' mailbox\n",box.name)
buf:=make([]byte,0,0x8000)


/*160:*/


//line amail.w:1958

c:=0
for len(v.msgs)> 0&&c<100{
msg:=v.msgs[0]
glog.V(debug).Infof("printing of '%s/%d' message with in the '%s' mailbox\n",msg.box.name,msg.id,box.name)
if box.threadMode(){


/*162:*/


//line amail.w:1986

{
for p:=msg.parent;p!=nil;p= p.parent{
buf= append(buf,levelmark...)
}
}



/*:162*/


//line amail.w:1964

}
c++


/*90:*/


//line amail.w:1208

glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n",msg.id,box.name)
buf= append(buf,fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n",
func()string{if msg.deleted{return deleted};return""}(),
func()string{if msg.box!=box{return fmt.Sprintf("%s/",msg.box.name)};return""}(),
msg.id,
func()string{if msg.unread{return newmark};return""}(),
msg.from,
msg.date,
msg.subject)...)




/*:90*/


//line amail.w:1967

v.msgs= v.msgs[1:]
if(v.flags&seek)==seek{
break
}
}
pcount+= c



/*:160*/


//line amail.w:1926

if _,err:=f.Write(buf);err!=nil{
glog.Errorf("can't write to 'data' file of the '%s' messgebox: %v\n",box.name,err)
}


/*159:*/


//line amail.w:1944

if!ontop{
glog.V(debug).Infof("pcount: %v, ontop: %v\n",pcount,ontop)


/*94:*/


//line amail.w:1249

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:94*/


//line amail.w:1947

if pcount>=100{
ontop= true
}
}




/*:159*/


//line amail.w:1930



/*158:*/


//line amail.w:1935

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:1939

}



/*:158*/


//line amail.w:1931

}



/*:157*/


//line amail.w:1739





/*:138*/



/*181:*/


//line amail.w:2261

case ids:=<-box.lch:
var msgs messages
for _,id:=range ids{
glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n",id,box.name)
p,ok:=box.all.Search(id)
if!ok{
glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n",id,box.name)
continue
}
msg:=box.all[p]
if msg.w==nil{
if msg.unread{


/*182:*/


//line amail.w:2289

msg.unread= false
box.unread.DeleteById(id)



/*:182*/


//line amail.w:2274



/*183:*/


//line amail.w:2295

if!box.thread&&box.shownew{


/*185:*/


//line amail.w:2318

box.eraseMessage(msg)




/*:185*/


//line amail.w:2297



/*184:*/


//line amail.w:2305

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*168:*/


//line amail.w:2050

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.deleted> 0{
return"UnDelmesg "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:168*/


//line amail.w:2308

w:=box.w
if box.deleted==0{


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2311

}else{


/*51:*/


//line amail.w:726

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2313

}
}



/*:184*/


//line amail.w:2298

}
msgs= append(msgs,msg)





/*:183*/


//line amail.w:2275



/*70:*/


//line amail.w:930

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:70*/


//line amail.w:2276

}
if err:=msg.open();err!=nil{
continue
}
}else{
glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n",id,box.name)
msg.w.WriteCtl("dot=addr\nshow")
}
}


/*148:*/


//line amail.w:1838

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:148*/


//line amail.w:2286




/*:181*/


//line amail.w:837

}
}
}



/*:61*/



/*69:*/


//line amail.w:924

func(box*mailbox)threadMode()bool{
return box.thread||box.showthreads&&!box.shownew
}



/*:69*/



/*122:*/


//line amail.w:1524

func(this idmessages)Search(messageid string)(int,bool){
pos:=sort.Search(len(this),func(i int)bool{return this[i].messageid<=messageid});
if pos!=len(this)&&this[pos].messageid==messageid{
return pos,true
}
return pos,false
}



/*:122*/



/*123:*/


//line amail.w:1534

func(this*idmessages)Insert(msg*message,pos int){
*this= append(*this,nil)
copy((*this)[pos+1:],(*this)[pos:])
(*this)[pos]= msg
}



/*:123*/



/*124:*/


//line amail.w:1544

func(this*idmessages)SearchInsert(msg*message)(int,bool){
pos,ok:=this.Search(msg.messageid)
if ok{
return pos,false
}
this.Insert(msg,pos)
return pos,true
}



/*:124*/



/*126:*/


//line amail.w:1558

func(this*idmessages)Delete(pos int)(*message,bool){
if pos<0||pos> len(*this)-1{
return nil,false
}
msg:=(*this)[pos]
*this= append((*this)[:pos],(*this)[pos+1:]...)
return msg,true
}



/*:126*/



/*129:*/


//line amail.w:1634

func(this idmessages)Len()int{
return len(this)
}

func(this idmessages)Less(i,j int)bool{
return this[i].date.Unix()<this[j].date.Unix()
}

func(this idmessages)Swap(i,j int){
t:=this[i]
this[i]= this[j]
this[j]= t
}




/*:129*/



/*142:*/


//line amail.w:1784

func getchildren(msg*message,dst messages,src messages)(messages,messages){


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:1786

for _,v:=range children{
dst= append(dst,v)
if p,ok:=src.Search(v.id);ok&&src[p]==v{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
dst,src= getchildren(v,dst,src)
}
return dst,src
}



/*:142*/



/*165:*/


//line amail.w:2023

func writeTag(w*goacme.Window,t string)error{
if w==nil{
return nil
}
tag,err:=w.File("tag")
if err!=nil{
return err
}
if err:=w.WriteCtl("cleartag");err!=nil{
return err
}
_,err= tag.Write([]byte(t))
return err
}



/*:165*/



/*170:*/


//line amail.w:2089

func clean(w*goacme.Window){
if err:=w.WriteAddr(wholefile);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",wholefile,err)
}else if data,err:=w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if _,err:=data.Write([]byte(""));err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
}
}



/*:170*/



/*177:*/


//line amail.w:2211

func(box*mailbox)search(str string)(msgs messages){
if len(str)==0{
return
}
f,err:=box.fid.Walk("search")
if err==nil{
err= f.Open(plan9.ORDWR)
}
if err!=nil{
glog.Errorf("can't open 'search' file: %s\n",err)
return
}
defer f.Close()
if _,err:=f.Write([]byte(str));err!=nil{
glog.Errorf("can't write to 'search' file: %s\n",err)
}
b:=bufio.NewReader(f)
for s,err:=b.ReadString(' ');err==nil||err==io.EOF;s,err= b.ReadString(' '){
s= strings.TrimSpace(s)
glog.V(debug).Infoln("search: ",s)
if num,err:=strconv.Atoi(s);err==nil{
if p,ok:=box.all.Search(num);ok{
msgs.Insert(box.all[p],0)
}
}
if err==io.EOF{
break
}
}
return
}



/*:177*/



/*186:*/


//line amail.w:2323

func(box*mailbox)eraseMessage(msg*message){
if box.w==nil{
return
}
glog.V(debug).Infof("removing the '%d' message of the '%s' mailbox from the '%s' mailbox\n",
msg.id,msg.box.name,box.name)


/*174:*/


//line amail.w:2147

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:174*/


//line amail.w:2330

if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file of the box '%s': %s\n",box.name,err)
}else if _,err:=data.Write([]byte{});err!=nil{
glog.Errorf("can't write to 'data' file of the box '%s': %s\n",box.name,err)
}
}



/*:186*/



/*189:*/


//line amail.w:2359

func split(s string)(strs[]string){
f:=strings.Fields(s)
m:=""
for _,v:=range f{
if strings.Contains(v,"@"){
m+= v
strs= append(strs,m)
m= ""
}else if v!="''"{
m+= v+" "
}
}
return
}



/*:189*/



/*190:*/


//line amail.w:2376

func(msg*message)open()(err error){
glog.V(debug).Infof("open: trying to open '%d' directory\n",msg.id)
bfid,err:=msg.box.fid.Walk(fmt.Sprintf("%d",msg.id))
if err!=nil{
glog.Errorf("can't walk to '%s/%d': %v\n",msg.box.name,msg.id,err)
return err
}
defer bfid.Close()
if msg.w==nil{
if msg.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
return err
}


/*198:*/


//line amail.w:2511

go func(){
glog.V(debug).Infof("starting a goroutine to process events from the '%d' message's window\n",msg.id)
for ev,err:=msg.w.ReadEvent();err==nil;ev,err= msg.w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
msg.w.UnreadEvent(ev)
continue
}
quote:=false
replyall:=false
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
msg.w.UnreadEvent(ev)
msg.w.Close()
msg.w= nil
return
case"Delmesg":
if!msg.deleted{
msg.deleted= true
msg.box.deleted++
msg.w.Del(true)
msg.w.Close()
msg.w= nil


/*146:*/


//line amail.w:1821

{
glog.V(debug).Infof("refresh a message '%d'\n",msg.id)
mrfch<-&refresh{seek,append(messages{},msg)}
}



/*:146*/


//line amail.w:2535

return
}
continue
case"UnDelmesg":
if msg.deleted{
msg.deleted= false
msg.box.deleted--


/*191:*/


//line amail.w:2408

msg.writeTag()



/*:191*/


//line amail.w:2543



/*146:*/


//line amail.w:1821

{
glog.V(debug).Infof("refresh a message '%d'\n",msg.id)
mrfch<-&refresh{seek,append(messages{},msg)}
}



/*:146*/


//line amail.w:2544

}
continue
case"Text":
if len(msg.text)!=0&&len(msg.html)!=0{
msg.showhtml= false
msg.open()
}
continue
case"Html":
if len(msg.text)!=0&&len(msg.html)!=0{
msg.showhtml= true
msg.open()
}
continue
case"Browser":


/*212:*/


//line amail.w:2866

{


/*217:*/


//line amail.w:2973



/*221:*/


//line amail.w:2994

once.Do(func(){

/*209:*/


//line amail.w:2810

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:209*/



/*218:*/


//line amail.w:2977

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:218*/



/*232:*/


//line amail.w:3283

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:232*/


//line amail.w:2995
})



/*:221*/


//line amail.w:2974




/*:217*/


//line amail.w:2868

dir:=fmt.Sprintf("%s/amail-%s/%s/%d",os.TempDir(),cuser,msg.box.name,msg.id)
if err:=os.MkdirAll(dir,0700);err!=nil{
glog.Errorf("can't create a directory '%s': %v\n",dir,err)
continue
}

if len(msg.files)==0{
if err:=saveFile(fmt.Sprintf("%s/%s/%d/%s",srv,msg.box.name,msg.id,msg.html),
fmt.Sprintf("%s/%d.html",dir,msg.id));err!=nil{
continue
}
}else{
if err:=msg.fixFile(dir);err!=nil{
continue
}
for _,v:=range msg.files{
saveFile(fmt.Sprintf("%s/%s/%d/%s/body",srv,msg.box.name,msg.id,v.path),
fmt.Sprintf("%s/%s",dir,v.name))
}

}

if p,err:=goplumb.Open("send",plan9.OWRITE);err!=nil{
glog.Errorf("can't open plumbing port 'send': %v\n",err)
}else if err:=p.SendText("amail","web",dir,fmt.Sprintf("file://%s/%d.html",dir,msg.id));err!=nil{
glog.Errorf("can't plumb a message '%s': %v\n",fmt.Sprintf("file://%s/%d.html",dir,msg.id),err)
}
}



/*:212*/


//line amail.w:2560

continue
case"Save":


/*222:*/


//line amail.w:2998

{
if len(ev.Arg)==0{
continue
}
f,err:=msg.box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl': %v\n",err)
continue
}
bs:=strings.Fields(ev.Arg)
for _,v:=range bs{
s:=fmt.Sprintf("save %s %d/",v,msg.id)
if _,err:=f.Write([]byte(s));err!=nil{
glog.Errorf("can't write '%s' to 'ctl': %v\n",s,err)
}
}
f.Close()
}




/*:222*/


//line amail.w:2563

continue
case"Q":
quote= true
fallthrough
case"Reply","Replyall":
if ev.Text=="Reply"{
args:=strings.Fields(ev.Arg)
for _,v:=range args{
if v=="all"{
replyall= true
}
}
}else if ev.Text=="Replyall"{
replyall= true
}


/*223:*/


//line amail.w:3023

{


/*224:*/


//line amail.w:3063

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*225:*/


//line amail.w:3075

go func(msg*message){
glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
for ev,err:=w.ReadEvent();err==nil;ev,err= w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
w.UnreadEvent(ev)
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
w.UnreadEvent(ev)
w.Close()
return
case"Post":


/*229:*/


//line amail.w:3162

{


/*231:*/


//line amail.w:3279



/*221:*/


//line amail.w:2994

once.Do(func(){

/*209:*/


//line amail.w:2810

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:209*/



/*218:*/


//line amail.w:2977

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:218*/



/*232:*/


//line amail.w:3283

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:232*/


//line amail.w:2995
})



/*:221*/


//line amail.w:3280




/*:231*/


//line amail.w:3164

w.Seek(0,0)
w.WriteAddr(wholefile)
ff,_:=w.File("xdata")
b:=bufio.NewReader(ff)
var to,cc,bcc,attach,include[]string
var subject string
for{
s,err:=b.ReadString('\n')
if err!=nil{
break
}
s= strings.TrimSpace(s)
if len(s)==0{
// an empty line, the rest is a body of the message
break
}
p:=strings.Index(s,":")
if p!=-1{
f:=strings.Split(s[p+1:],",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
switch strings.ToLower(s[:p]){
case"from","to":
to= append(to,f...)
case"cc":
cc= append(cc,f...)
case"bcc":
bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%q",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
to= append(to,f...)
}
}
args:=append([]string{},"-8")
if msg!=nil{
args= append(args,"-R",fmt.Sprintf("%s/%d",msg.box.name,msg.id))
}
if len(subject)!=0{
args= append(args,"-s",subject)
}
for _,v:=range include{
args= append(args,"-A",v)
}
for _,v:=range attach{
args= append(args,"-a",v)
}

c:=exec.Command(plan9dir+"/bin/upas/marshal",args...)
p,err:=c.StdinPipe()
if err!=nil{
glog.Errorf("can't get a stdin pipe: %v\n",err)
continue
}
if err:=c.Start();err!=nil{
glog.Errorf("can't start 'upas/marshal': %v\n",err)
continue
}
if len(to)!=0{
if _,err:=fmt.Fprintln(p,"To:",strings.Join(to,","));err!=nil{
glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("to is written")
if len(cc)!=0{
if _,err:=fmt.Fprintln(p,"CC:",strings.Join(cc,","));err!=nil{
glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("cc is written")
if len(bcc)!=0{
if _,err:=fmt.Fprintln(p,"BCC:",strings.Join(bcc,","));err!=nil{
glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("bcc is written")
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("writing '%s':%v",s,err)

p.Write([]byte(s))
if err==io.EOF{
p.Write([]byte("\n"))
break
}
}
glog.V(debug).Infoln("body is written")
p.Write([]byte("\n"))
p.Close()
c.Wait()
w.Del(true)
w.Close()
}





/*:229*/


//line amail.w:3090

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:225*/


//line amail.w:3072




/*:224*/


//line amail.w:3025

name:=fmt.Sprintf("Amail/%s/%d/%sReply%s",
msg.box.name,
msg.id,
func()string{if quote{return"Q"};return""}(),
func()string{if replyall{return"all"};return""}())


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:3031

buf:=make([]byte,0,0x8000)
buf= append(buf,fmt.Sprintf("To: %s\n",msg.from)...)
if replyall{
for _,v:=range msg.to{
buf= append(buf,fmt.Sprintf("To: %s\n",v)...)
}
for _,v:=range msg.cc{
buf= append(buf,fmt.Sprintf("To: %s\n",v)...)
}
}
buf= append(buf,fmt.Sprintf("Subject: %s%s\n",
func()string{
if!strings.Contains(msg.subject,"Re:"){
return"Re: "
}
return""
}(),
msg.subject)...)
if quote{
buf= append(buf,'\n')


/*226:*/


//line amail.w:3098

if len(msg.text)!=0{
fn:=fmt.Sprintf("%d/%s",msg.id,msg.text)
f,err:=msg.box.fid.Walk(fn)
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open '%s/%s/%s': %v\n",srv,msg.box.name,fn)
continue
}


/*227:*/


//line amail.w:3116

{
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
buf= append(buf,'>',' ')
buf= append(buf,s...)
if err==io.EOF{
break
}
}
}



/*:227*/


//line amail.w:3109

f.Close()
}else if len(msg.html)!=0{


/*228:*/


//line amail.w:3129

{
c1:=exec.Command("9p","read",fmt.Sprintf("%s/%s/%d/%s",srv,msg.box.name,msg.id,msg.html))
c2:=exec.Command("htmlfmt","-cutf-8")
f,err:=c2.StdoutPipe()
if err!=nil{
glog.Errorf("can't get a stdout pipe: %v\n",err)
}
c2.Stdin,err= c1.StdoutPipe()
if err!=nil{
glog.Errorf("can't get a stdout pipe: %v\n",err)
f.(io.Closer).Close()
continue
}
if err= c2.Start();err!=nil{
glog.Errorf("can't start 'htmlfmt': %v\n",err)
f.(io.Closer).Close()
continue
}
if err= c1.Start();err!=nil{
glog.Errorf("can't run '9p': %v\n",err)
c2.Wait()
f.(io.Closer).Close()
continue
}


/*227:*/


//line amail.w:3116

{
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
buf= append(buf,'>',' ')
buf= append(buf,s...)
if err==io.EOF{
break
}
}
}



/*:227*/


//line amail.w:3154

c1.Wait()
c2.Wait()
f.(io.Closer).Close()
}



/*:228*/


//line amail.w:3112

}



/*:226*/


//line amail.w:3052

}else{
buf= append(buf,fmt.Sprintf("Include: Mail/%s/%d/raw\n",msg.box.name,msg.id)...)

}
buf= append(buf,'\n')
w.Write(buf)


/*94:*/


//line amail.w:1249

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:94*/


//line amail.w:3059

}



/*:223*/


//line amail.w:2579

continue
case"Up":
if msg.parent!=nil{


/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:2583

name:=msg.parent.box.name
id:=msg.parent.id


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:2586



/*87:*/


//line amail.w:1177

lch<-&msgs



/*:87*/


//line amail.w:2587

}
continue
case"Down":


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:2591

if len(children)!=0{


/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:2593

name:=children[0].box.name
id:=children[0].id


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:2596



/*87:*/


//line amail.w:1177

lch<-&msgs



/*:87*/


//line amail.w:2597

}
continue
case"Prev":


/*192:*/


//line amail.w:2412

pmsg:=msg.prev()



/*:192*/


//line amail.w:2601

if pmsg!=nil{


/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:2603

name:=pmsg.box.name
id:=pmsg.id


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:2606



/*87:*/


//line amail.w:1177

lch<-&msgs



/*:87*/


//line amail.w:2607

}
continue
case"Next":


/*194:*/


//line amail.w:2433

nmsg:=msg.next()



/*:194*/


//line amail.w:2611

if nmsg!=nil{


/*84:*/


//line amail.w:1157

msgs:=make(msgmap)



/*:84*/


//line amail.w:2613

name:=nmsg.box.name
id:=nmsg.id


/*85:*/


//line amail.w:1161

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:85*/


//line amail.w:2616



/*87:*/


//line amail.w:1177

lch<-&msgs



/*:87*/


//line amail.w:2617

}
continue
}
}else if(ev.Type&goacme.Look)==goacme.Look{
}
msg.w.UnreadEvent(ev)

}
}()



/*:198*/


//line amail.w:2390

}else{


/*199:*/


//line amail.w:2629

glog.V(debug).Infof("clean the '%s/%d' message's window\n",msg.box.name,msg.id)
clean(msg.w)



/*:199*/


//line amail.w:2392

}
buf:=make([]byte,0,0x8000)


/*197:*/


//line amail.w:2501

{
glog.V(debug).Infof("composing a header of the '%d' message\n",msg.id)
buf= append(buf,fmt.Sprintf("From: %s\nDate: %s\nTo: %s\n%sSubject: %s\n\n\n",
msg.from,msg.date,strings.Join(msg.to,", "),
func()string{if len(msg.cc)!=0{return fmt.Sprintf("CC: %s\n",strings.Join(msg.cc,", "))};return""}(),
msg.subject)...)
}



/*:197*/


//line amail.w:2395



/*204:*/


//line amail.w:2670

{
if len(msg.text)==0&&len(msg.html)==0{
if err= msg.bodyPath(bfid,"");err!=nil{
glog.Errorf("can't ged a body path of '%d': %v\n",msg.id,err)
}
glog.V(debug).Infof("paths for bodies of the '%d' message have been found: text-'%s', html-'%s'\n",
msg.id,msg.text,msg.html)

}
if len(msg.text)!=0&&!msg.showhtml{
glog.V(debug).Infof("using a path for a text body of the '%d' message: '%s'\n",msg.id,msg.text)
if buf,err= readAll(bfid,msg.text,buf);err!=nil{
glog.Errorf("can't read '%s': %v\n",msg.text,err)
return
}
}else if len(msg.html)!=0{
glog.V(debug).Infof("using a path for a html body of the '%d' message: '%s'\n",msg.id,msg.html)
msg.w.Write(buf)
buf= nil
c1:=exec.Command("9p","read",fmt.Sprintf("%s/%s/%d/%s",srv,msg.box.name,msg.id,msg.html))
c2:=exec.Command("htmlfmt","-cutf-8")
c2.Stdout,_= msg.w.File("body")
c2.Stdin,err= c1.StdoutPipe()
if err!=nil{
glog.Errorf("can't get a stdout pipe: %v\n",err)
return
}
if err= c2.Start();err!=nil{
glog.Errorf("can't start 'htmlfmt': %v\n",err)
return
}
if err= c1.Run();err!=nil{
glog.Errorf("can't run '9p': %v\n",err)
c2.Wait()
return
}
if err= c2.Wait();err!=nil{
glog.Errorf("can't wait of completion 'htmlfmt': %v\n",err)
return
}
}


/*208:*/


//line amail.w:2806



/*221:*/


//line amail.w:2994

once.Do(func(){

/*209:*/


//line amail.w:2810

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:209*/



/*218:*/


//line amail.w:2977

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:218*/



/*232:*/


//line amail.w:3283

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:232*/


//line amail.w:2995
})



/*:221*/


//line amail.w:2807




/*:208*/


//line amail.w:2712

for _,v:=range msg.files{
buf= append(buf,fmt.Sprintf("\n===> %s (%s)\n",v.path,v.mimetype)...)
buf= append(buf,fmt.Sprintf("\t9p read %s/%s/%d/%sbody > '%s/%s'\n",srv,msg.box.name,msg.id,v.path,home,v.name)...)
}
}



/*:204*/


//line amail.w:2396

w:=msg.w
name:=fmt.Sprintf("Amail/%s/%d",msg.box.name,msg.id)


/*53:*/


//line amail.w:747

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2399



/*191:*/


//line amail.w:2408

msg.writeTag()



/*:191*/


//line amail.w:2400

w.Write(buf)


/*50:*/


//line amail.w:717

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2402



/*94:*/


//line amail.w:1249

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:94*/


//line amail.w:2403

return
}



/*:190*/



/*193:*/


//line amail.w:2416

func(this*message)prev()(pmsg*message){
if this.parent==nil{
return
}
msg:=this.parent


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:2422

for _,v:=range children{
if v==this{
break
}
pmsg= v
}
return
}



/*:193*/



/*195:*/


//line amail.w:2437

func(this*message)next()(nmsg*message){
if this.parent==nil{
return
}
msg:=this.parent


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:2443

for i:=0;i<len(children);i++{
if children[i]!=this{
continue
}
i++
if i<len(children){
nmsg= children[i]
}
break
}
return
}



/*:195*/



/*196:*/


//line amail.w:2458

func(msg*message)writeTag(){
glog.V(debug).Infof("writing a tag of the '%d' message's window\n",msg.id)
if err:=writeTag(msg.w,fmt.Sprintf(" Q Reply all %s %s%s%s%s%s%sSave ",
func()string{if msg.deleted{return"UnDelmesg"}else{return"Delmesg"}}(),
func()string{
if len(msg.text)==0||len(msg.html)==0{
return""
}else if msg.showhtml{
return"Text "
}else{
return"Html "
}
}(),
func()string{if len(msg.html)!=0{return"Browser "};return""}(),
func()string{if msg.parent!=nil{return"Up "};return""}(),
func()string{


/*131:*/


//line amail.w:1666

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:131*/


//line amail.w:2475

if len(children)!=0{
return"Down "
}
return""
}(),
func()string{


/*192:*/


//line amail.w:2412

pmsg:=msg.prev()



/*:192*/


//line amail.w:2482

if pmsg!=nil{
return"Prev "
}
return""
}(),
func()string{


/*194:*/


//line amail.w:2433

nmsg:=msg.next()



/*:194*/


//line amail.w:2489

if nmsg!=nil{
return"Next "
}
return""
}()));
err!=nil{
glog.Errorf("can't set a tag of the message window: %v",err)
}
}



/*:196*/



/*205:*/


//line amail.w:2720

func(msg*message)bodyPath(bfid*client.Fid,path string)error{
glog.V(debug).Infof("getting a path for a body of the '%d' message\n",msg.id)
t,err:=readString(bfid,path+"type")
if err!=nil{
return err
}
switch t{
case"message/rfc822",
"text",
"text/plain",
"text/richtext",
"text/tab-separated-values":
if len(msg.text)==0{
msg.text= path+"body"
glog.V(debug).Infof("a path for a text body of the '%d' message: '%s'\n",msg.id,t)
return nil
}
case"text/html":
if len(msg.html)==0{
msg.html= path+"body"
glog.V(debug).Infof("a path for a html body of the '%d' message: '%s'\n",msg.id,t)
return nil
}
case"multipart/mixed",
"multipart/alternative",
"multipart/related",
"multipart/signed",
"multipart/report":
for c:=1;;c++{
if err= msg.bodyPath(bfid,fmt.Sprintf("%s%d/",path,c));err!=nil{
break
}
}
return nil
}
glog.V(debug).Infof("trying to read '%d/%sfilename'\n",msg.id,path)
if n,err:=readString(bfid,path+"filename");err==nil{
f:=&file{name:n,mimetype:t,path:path,}
if len(n)==0{
f.name= "attachment"
}else if cid,ok:=msg.getCID(path);ok{
msg.cids[cid]= f
}
msg.files= append(msg.files,f)
}
return nil
}



/*:205*/



/*206:*/


//line amail.w:2770

func(msg*message)getCID(path string)(string,bool){
src:=fmt.Sprintf("%d/%smimeheader",msg.id,path)
glog.V(debug).Infof("getting of cids for path '%s'\n",src)
fid,err:=msg.box.fid.Walk(src)
if err==nil{
err= fid.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open '%s': %v\n",src,err)
return"",false
}
defer fid.Close()
fid.Seek(0,0)
b:=bufio.NewReader(fid)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("looking for a cid in '%s'\n",s)
if strings.HasPrefix(s,"Content-ID: <"){
s= s[len("Content-ID: <"):len(s)-2]
glog.V(debug).Infof("found a cid '%s'\n",s)
return s,true
}
if err==io.EOF{
break
}
}
return"",false
}





/*:206*/



/*210:*/


//line amail.w:2820

func readString(pfid*client.Fid,name string)(str string,err error){
glog.V(debug).Infof("readString: trying to open '%s'\n",name)
f,err:=pfid.Walk(name)
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
return
}
defer f.Close()
str,err= bufio.NewReader(f).ReadString('\n')
if err!=nil&&err!=io.EOF{
return
}
return str,nil
}



/*:210*/



/*211:*/


//line amail.w:2839

func readAll(pfid*client.Fid,name string,buf[]byte)([]byte,error){
glog.V(debug).Infof("readAll: trying to open '%s'\n",name)
f,err:=pfid.Walk(name)
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
return buf,err
}
defer f.Close()
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil;s,err= b.ReadString('\n'){
if strings.HasSuffix(s,"\r\n"){
s= strings.TrimRight(s,"\r\n")
s+= "\n"
}
buf= append(buf,s...)
}
return buf,nil
}




/*:211*/



/*213:*/


//line amail.w:2899

func saveFile(src,dst string)error{
var err error
c:=exec.Command("9p","read",src)
f,err:=os.OpenFile(dst,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0600)
if err!=nil{
glog.Errorf("can't create a file '%s': %v\n",dst,err)
return err
}
defer f.Close()
c.Stdout= f
if err= c.Run();err!=nil{
glog.Errorf("can't run '9p': %v\n",err)
}
return err
}



/*:213*/



/*214:*/


//line amail.w:2917

func(msg*message)fixFile(dir string)error{
src:=fmt.Sprintf("%d/%s",msg.id,msg.html)
dst:=fmt.Sprintf("%s/%d.html",dir,msg.id)
df,err:=os.OpenFile(dst,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0600)
if err!=nil{
glog.Errorf("can't create a file '%s': %v\n",dst,err)
return err
}
defer df.Close()
fid,err:=msg.box.fid.Walk(src)
if err==nil{
err= fid.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open to '%s': %v\n",src,err)
return err
}
defer fid.Close()
b:=bufio.NewReader(fid)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
p:=0
for b:=strings.Index(s[p:],"\"cid:");b!=-1;b= strings.Index(s[p:],"\"cid:"){
b+= p
e:=strings.Index(s[b+1:],"\"")
if e==-1{
break
}
e++
glog.V(debug).Infof("len(s): %v, p: %v, b: %v, e: %v\n",len(s),p,b,e)
cid:=s[b+5:b+e]
glog.V(debug).Infof("cid: %s\n",cid)
if f,ok:=msg.cids[cid];ok{
glog.V(debug).Infof("found a cid: %s, replace '%s' by '%s'\n",cid,s[b+1:b+e],f.name)
s= strings.Replace(s,s[b+1:b+e],f.name,1)
}else{
p= b+e
}
}
df.Write([]byte(s))
if err==io.EOF{
break
}
}
return err
}



/*:214*/


