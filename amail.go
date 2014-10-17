

/*2:*/


//line amail.w:32

//line license:1

// This file is part of Amail version 0.92
// Author Alexander Sychev
//
// Copyright (c) 2013, 2014 Alexander Sychev. All rights reserved.
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
//line amail.w:34

package main

import(


/*11:*/


//line amail.w:132

"flag"
"fmt"
"os"
"strings"
"sort"



/*:11*/



/*13:*/


//line amail.w:166

"unicode"
"unicode/utf8"



/*:13*/



/*15:*/


//line amail.w:184

"bitbucket.org/santucco/goplan9-clone/plan9/client"
"github.com/golang/glog"



/*:15*/



/*28:*/


//line amail.w:312

"io"
"bufio"



/*:28*/



/*31:*/


//line amail.w:358

"bitbucket.org/santucco/goplumb"
"bitbucket.org/santucco/goplan9-clone/plan9"



/*:31*/



/*39:*/


//line amail.w:543

"bitbucket.org/santucco/goacme"



/*:39*/



/*46:*/


//line amail.w:692

"strconv"



/*:46*/



/*92:*/


//line amail.w:1290

"time"



/*:92*/



/*121:*/


//line amail.w:1603

"errors"



/*:121*/



/*208:*/


//line amail.w:2916

"os/exec"



/*:208*/



/*223:*/


//line amail.w:3250

"os/user"



/*:223*/



/*227:*/


//line amail.w:3271

"sync"



/*:227*/


//line amail.w:38

)

type(


/*19:*/


//line amail.w:208

mailbox struct{
name string


/*21:*/


//line amail.w:234

all messages
unread messages
mch chan int
dch chan int



/*:21*/



/*47:*/


//line amail.w:696

fid*client.Fid
total int



/*:47*/



/*71:*/


//line amail.w:991

shownew bool
showthreads bool
ech<-chan*goacme.Event
w*goacme.Window
cch chan bool



/*:71*/



/*76:*/


//line amail.w:1042

thread bool



/*:76*/



/*106:*/


//line amail.w:1457

mdch chan messages



/*:106*/



/*140:*/


//line amail.w:1900

rfch chan*refresh
irfch chan*refresh
reset bool



/*:140*/



/*158:*/


//line amail.w:2083

pos int



/*:158*/



/*173:*/


//line amail.w:2260

deleted int



/*:173*/



/*186:*/


//line amail.w:2507

ach chan*struct{ids[]int;a action}



/*:186*/


//line amail.w:211

}

mailboxes[]*mailbox

message struct{
id int


/*29:*/


//line amail.w:317

unread bool
box*mailbox



/*:29*/



/*63:*/


//line amail.w:923

deleted bool



/*:63*/



/*93:*/


//line amail.w:1294

from string
date time.Time
subject string



/*:93*/



/*120:*/


//line amail.w:1598

inreplyto string
messageid string



/*:120*/



/*188:*/


//line amail.w:2515

w*goacme.Window



/*:188*/



/*195:*/


//line amail.w:2605

to[]string
cc[]string



/*:195*/



/*210:*/


//line amail.w:2938

text string
html string
showhtml bool
files[]*file
cids map[string]*file



/*:210*/


//line amail.w:218

}

messages[]*message




/*:19*/



/*81:*/


//line amail.w:1205

msgmap map[string][]int



/*:81*/



/*82:*/


//line amail.w:1211

action int



/*:82*/



/*117:*/


//line amail.w:1578

idmessages[]*message
rootmsg*message
parentmsg*message



/*:117*/



/*118:*/


//line amail.w:1585

idlinks struct{
msgs messages
parent*idlinks
children[]*idlinks
}



/*:118*/



/*138:*/


//line amail.w:1881

refreshFlags int

refresh struct{
flags refreshFlags
msgs messages
}



/*:138*/



/*209:*/


//line amail.w:2920

file struct{
name string
mimetype string
path string
}



/*:209*/


//line amail.w:42

)



/*42:*/


//line amail.w:572

const mailboxfmt= "%-30s\t%10d\t%10d\n"
const mailboxfmtprc= "%-30s\t%10d\t%10d\t%d%%\n"
const wholefile= "0,$"



/*:42*/



/*83:*/


//line amail.w:1215

const(
view action= iota


/*99:*/


//line amail.w:1350

del
undel



/*:99*/



/*113:*/


//line amail.w:1507

seen



/*:113*/


//line amail.w:1218

)



/*:83*/



/*139:*/


//line amail.w:1892

const(
seek refreshFlags= 1<<iota
insert refreshFlags= 1<<iota
exact refreshFlags= 1<<iota
)



/*:139*/



/*163:*/


//line amail.w:2121

const eof= "$"



/*:163*/



/*179:*/


//line amail.w:2357

const bof= "#0-"
const eol= "+#0"



/*:179*/


//line amail.w:45


var(


/*3:*/


//line amail.w:81

exit chan bool= make(chan bool)



/*:3*/



/*6:*/


//line amail.w:99

wch chan int= make(chan int,100)
wcount int



/*:6*/



/*10:*/


//line amail.w:124

shownew bool
showthreads bool
levelmark string
newmark string
skipboxes[]string



/*:10*/



/*16:*/


//line amail.w:189

fsys*client.Fsys
rfid*client.Fid
srv string= "mail"



/*:16*/



/*20:*/


//line amail.w:225

boxes mailboxes



/*:20*/



/*23:*/


//line amail.w:252

mch= make(chan*struct{name string;id int},100)
dch= make(chan*struct{name string;id int},100)
bch= make(chan string,10)
rfch= make(chan*mailbox,100)



/*:23*/



/*40:*/


//line amail.w:547

mw*goacme.Window
ech<-chan*goacme.Event



/*:40*/



/*84:*/


//line amail.w:1222

ach= make(chan*struct{m msgmap;a action},100)



/*:84*/



/*90:*/


//line amail.w:1273

deleted= "(deleted)-"



/*:90*/



/*109:*/


//line amail.w:1483

mdch chan messages= make(chan messages,100)



/*:109*/



/*119:*/


//line amail.w:1593

idmap= make(map[string]*idlinks)
idch= make(chan struct{msg*message;val interface{}},100)



/*:119*/



/*156:*/


//line amail.w:2069

mrfch chan*refresh= make(chan*refresh)



/*:156*/



/*215:*/


//line amail.w:3084

home string



/*:215*/



/*224:*/


//line amail.w:3254

cuser string



/*:224*/



/*228:*/


//line amail.w:3275

once sync.Once



/*:228*/



/*239:*/


//line amail.w:3575

plan9dir string



/*:239*/


//line amail.w:48

debug glog.Level= 1
)

func main(){
glog.V(debug).Infoln("main")
defer glog.V(debug).Infoln("main is done")


/*12:*/


//line amail.w:140

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


//line amail.w:172

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


//line amail.w:156

if len(skip)> 0{
skipboxes= strings.Split(skip,", ")
sort.Strings(skipboxes)
glog.V(debug).Infof("these mailboxes will be skipped: %v\n",skipboxes)

}
}



/*:12*/


//line amail.w:55



/*17:*/


//line amail.w:195

{
glog.V(debug).Infoln("try to open mailfs")
var err error
if fsys,err= client.MountService(srv);err!=nil{
glog.Errorf("can't mount mailfs: %v\n",err)
os.Exit(1)
}
}




/*:17*/


//line amail.w:56



/*32:*/


//line amail.w:365

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


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:380

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
s:=strings.TrimLeft(string(m.Data),"Mail/")
n:=strings.LastIndex(s,"/")
if n==-1{
glog.Warning("can't found a number of message in '%s'\n",s)
continue
}
num,err:=strconv.Atoi(s[n+1:])
if err!=nil{
glog.Error(err)
continue
}
if v=="new"{
glog.V(debug).Infof("'%d' is a new message in the '%s' mailbox\n",num,s[:n])
mch<-&struct{name string;id int}{name:s[:n],id:num}
}else if v=="delete"{
glog.V(debug).Infof("'%d' is a deleted message in the '%s' mailbox\n",num,s[:n])
dch<-&struct{name string;id int}{name:s[:n],id:num}
}
}
}
}()
}
}
}



/*:32*/



/*33:*/


//line amail.w:423

{
glog.V(debug).Infoln("trying to open 'sendmail' plumbing port")
if sm,err:=goplumb.Open("sendmail",plan9.OREAD);err!=nil{
glog.V(debug).Infof("can't open plumb/sendmail: %s\n",err)
}else{
sch,err:=sm.MessageChannel(0)
if err!=nil{
glog.Errorf("can't get message channal for plumb/sendmail: %s\n",err)
}else{
go func(){
defer sm.Close()
defer glog.V(debug).Infoln("plumbing goroutine is done")
for{
select{


/*4:*/


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:438

case m,ok:=<-sch:
if!ok{
glog.Warningln("it seems plumber has finished")
sch= nil
return
}
glog.V(debug).Infof("a plumbing message has been received: %v\n",m)
var msg*message


/*232:*/


//line amail.w:3348

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*233:*/


//line amail.w:3360

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


/*237:*/


//line amail.w:3451

{


/*240:*/


//line amail.w:3579



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3580




/*:240*/


//line amail.w:3453

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
case"to":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3478

to= append(to,f...)
case"cc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3481

cc= append(cc,f...)
case"bcc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3484

bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%s",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3499

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



/*:237*/


//line amail.w:3375

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:233*/


//line amail.w:3357




/*:232*/


//line amail.w:447

name:=fmt.Sprintf("Amail/New")


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:449

addr:=fmt.Sprintf("To: %s\n\n",string(m.Data))
w.Write([]byte(addr))


/*244:*/


//line amail.w:3602

writeSignature(w,nil)




/*:244*/


//line amail.w:452

}
}
}()
}
}
}



/*:33*/


//line amail.w:57



/*25:*/


//line amail.w:276

glog.V(debug).Infoln("initialization of root of mailfs")
var err error
rfid,err= fsys.Walk(".")
if err!=nil{
glog.Errorf("can't open mailfs: %v\n",err)
os.Exit(1)
}
defer rfid.Close()




/*:25*/


//line amail.w:58



/*123:*/


//line amail.w:1626

go func(){
for{
select{


/*4:*/


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:1630

case v:=<-idch:
if v.val==nil{


/*125:*/


//line amail.w:1688

{
val,ok:=idmap[v.msg.messageid]
if!ok{
continue
}
for i,_:=range val.msgs{
if val.msgs[i]==v.msg{
val.msgs.Delete(i)
break
}
}
if len(val.msgs)> 0{
continue
}
if val.parent!=nil{
for i,_:=range val.parent.children{
if val.parent.children[i]==val{
val.parent.children= append(val.parent.children[:i],val.parent.children[i+1:]...)
break
}
}
}
for _,ch:=range val.children{
ch.parent= nil
}
if len(val.children)==0{
delete(idmap,v.msg.messageid)
}
}



/*:125*/


//line amail.w:1633

}else if ch,ok:=v.val.(chan bool);ok{


/*124:*/


//line amail.w:1656

{
glog.V(debug).Infof("appending a '%s' ('%s/%d') message to idmap\n",v.msg.messageid,v.msg.box.name,v.msg.id)
val,ok:=idmap[v.msg.messageid]
if!ok{
glog.V(debug).Infof("'%s' ('%s/%d') message  doesn't exist, creating\n",v.msg.messageid,v.msg.box.name,v.msg.id)
val= new(idlinks)
idmap[v.msg.messageid]= val
}
if len(val.msgs)> 0{
glog.V(debug).Infof("%v(%v) is a duplicate of %v(%v)\n",
v.msg.id,v.msg.messageid,
val.msgs[0].id,val.msgs[0].messageid)
}
val.msgs= append(val.msgs,v.msg)

if len(v.msg.inreplyto)!=0&&len(val.msgs)==1{
pval,ok:=idmap[v.msg.inreplyto]
if!ok{
pval= new(idlinks)
idmap[v.msg.inreplyto]= pval
}

pval.children= append(pval.children,val)
val.parent= pval
}
ch<-true
}



/*:124*/


//line amail.w:1635

}else if ch,ok:=v.val.(chan idmessages);ok{


/*127:*/


//line amail.w:1737

{
if m,ok:=idmap[v.msg.messageid];ok{
var children idmessages
for _,val:=range m.children{


/*130:*/


//line amail.w:1786

var msg*message
if val!=nil&&len(val.msgs)> 0{
msg= val.msgs[0]
for i,_:=range val.msgs{
if val.msgs[i].box==v.msg.box{
msg= val.msgs[i]
break
}
}
}



/*:130*/


//line amail.w:1742

if msg!=nil{
children= append(children,msg)
}
}
sort.Sort(children)
glog.V(debug).Infof("sending %d children for '%s'\n",len(children),v.msg.messageid)
ch<-children
}else{
glog.V(debug).Infof("'%s' is not found\n",v.msg.messageid)
ch<-nil
}
}



/*:127*/


//line amail.w:1637

}else if ch,ok:=v.val.(chan parentmsg);ok{


/*129:*/


//line amail.w:1767

{
if val,ok:=idmap[v.msg.messageid];!ok{
glog.V(debug).Infof("'%s' is not found\n",v.msg.messageid)
ch<-nil
}else if val.parent==nil||len(val.parent.msgs)==0{
glog.V(debug).Infof("'%s' hasn't got a parent\n",v.msg.messageid)
ch<-nil
}else{
val= val.parent


/*130:*/


//line amail.w:1786

var msg*message
if val!=nil&&len(val.msgs)> 0{
msg= val.msgs[0]
for i,_:=range val.msgs{
if val.msgs[i].box==v.msg.box{
msg= val.msgs[i]
break
}
}
}



/*:130*/


//line amail.w:1777

if msg!=nil{
glog.V(debug).Infof("sending parent '%s' for '%s'\n",msg.messageid,v.msg.messageid)
}
ch<-msg
}
}



/*:129*/


//line amail.w:1639

}else if ch,ok:=v.val.(chan rootmsg);ok{


/*132:*/


//line amail.w:1809

{
if val,ok:=idmap[v.msg.messageid];ok{
for val.parent!=nil&&len(val.parent.msgs)> 0{
val= val.parent
}


/*130:*/


//line amail.w:1786

var msg*message
if val!=nil&&len(val.msgs)> 0{
msg= val.msgs[0]
for i,_:=range val.msgs{
if val.msgs[i].box==v.msg.box{
msg= val.msgs[i]
break
}
}
}



/*:130*/


//line amail.w:1815

if msg==nil{
msg= v.msg
}
ch<-rootmsg(msg)
}else{
glog.V(debug).Infof("'%s' is not found\n",v.msg.messageid)
ch<-nil
}
}



/*:132*/


//line amail.w:1641

}else if ch,ok:=v.val.(chan int);ok{


/*134:*/


//line amail.w:1837

{
if val,ok:=idmap[v.msg.messageid];ok{
level:=0
for val.parent!=nil&&len(val.parent.msgs)> 0{
val= val.parent
level++
}
glog.V(debug).Infof("sending level '%d' for '%s' ('%s/%d')\n",level,v.msg.messageid,v.msg.box.name,v.msg.id)
ch<-level
}else{
glog.V(debug).Infof("'%s' is not found\n",v.msg.messageid)
ch<-0
}
}



/*:134*/


//line amail.w:1643

}
}
}
}()



/*:123*/


//line amail.w:59

if len(flag.Args())> 0{


/*35:*/


//line amail.w:494

go func(){
glog.V(debug).Infoln("start a main message loop")
defer glog.V(debug).Infoln("main message loop is done")
for{
select{


/*4:*/


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:500

case b:=<-rfch:
if b==nil{


/*43:*/


//line amail.w:578

if mw!=nil{
glog.V(debug).Infoln("printing of the mailboxes")
if err:=mw.WriteAddr(wholefile);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",wholefile,err)
}else if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else{
for _,v:=range boxes{
if v.total==len(v.all){
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,len(v.unread),len(v.all))))
}else if v.total!=0&&len(v.all)*100/v.total> 0{
data.Write([]byte(fmt.Sprintf(mailboxfmtprc,v.name,len(v.unread),len(v.all),len(v.all)*100/v.total)))
}else{
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,0,0)))
}
}
}
w:=mw


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:597



/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:598

}




/*:43*/


//line amail.w:503

}else{


/*45:*/


//line amail.w:659

glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n",b.name,len(b.all),b.total)
if mw!=nil{
if len(b.all)!=b.total&&b.total/100!=0&&len(b.all)%(b.total/100)!=0{
continue
}

if err:=mw.WriteAddr("0/^%s.*\\n/",escape(b.name));err!=nil{
glog.V(debug).Infof("can't write to 'addr' file: %s\n",err)
continue
}

if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if len(b.all)==b.total{
if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmt,b.name,len(b.unread),len(b.all))));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
w:=mw


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:680



/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:681

}else if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,
b.name,
len(b.unread),
len(b.all),
len(b.all)*100/b.total)));err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
}
}



/*:45*/


//line amail.w:505

}


/*44:*/


//line amail.w:604

case ev,ok:=<-ech:
glog.V(debug).Infof("an event from main window has been received: %v\n",ev)
if!ok{
ech= nil


/*9:*/


//line amail.w:118

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:609

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


//line amail.w:91

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:626

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


/*171:*/


//line amail.w:2233

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:171*/


//line amail.w:638

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


/*73:*/


//line amail.w:1005

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:73*/


//line amail.w:649

continue
}
}
mw.UnreadEvent(ev)




/*:44*/


//line amail.w:507

}
}
}()



/*:35*/


//line amail.w:61

for _,name:=range flag.Args(){


/*36:*/


//line amail.w:513

glog.V(debug).Infof("creating a '%s' mailbox\n",name)
box:=&mailbox{name:name,

/*22:*/


//line amail.w:241

mch:make(chan int,100),
dch:make(chan int,100),




/*:22*/



/*72:*/


//line amail.w:999

shownew:shownew,
showthreads:showthreads,
cch:make(chan bool,100),



/*:72*/



/*107:*/


//line amail.w:1461

mdch:make(chan messages,100),



/*:107*/



/*141:*/


//line amail.w:1906

rfch:make(chan*refresh,100),
irfch:make(chan*refresh,100),



/*:141*/



/*187:*/


//line amail.w:2511

ach:make(chan*struct{ids[]int;a action},100),



/*:187*/


//line amail.w:515
}
boxes= append(boxes,box)
sort.Sort(boxes)



/*:36*/


//line amail.w:63



/*59:*/


//line amail.w:860

go box.loop()



/*:59*/


//line amail.w:64



/*73:*/


//line amail.w:1005

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:73*/


//line amail.w:65

}
}else{


/*41:*/


//line amail.w:552

glog.V(debug).Infoln("creating the main window")
defer goacme.DeleteAll()

var err error
if mw,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
name:="Amail"
w:=mw


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:563

if ech,err= mw.EventChannel(0,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*171:*/


//line amail.w:2233

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:171*/


//line amail.w:568



/*8:*/


//line amail.w:113

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:569




/*:41*/


//line amail.w:68



/*35:*/


//line amail.w:494

go func(){
glog.V(debug).Infoln("start a main message loop")
defer glog.V(debug).Infoln("main message loop is done")
for{
select{


/*4:*/


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:500

case b:=<-rfch:
if b==nil{


/*43:*/


//line amail.w:578

if mw!=nil{
glog.V(debug).Infoln("printing of the mailboxes")
if err:=mw.WriteAddr(wholefile);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",wholefile,err)
}else if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else{
for _,v:=range boxes{
if v.total==len(v.all){
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,len(v.unread),len(v.all))))
}else if v.total!=0&&len(v.all)*100/v.total> 0{
data.Write([]byte(fmt.Sprintf(mailboxfmtprc,v.name,len(v.unread),len(v.all),len(v.all)*100/v.total)))
}else{
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,0,0)))
}
}
}
w:=mw


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:597



/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:598

}




/*:43*/


//line amail.w:503

}else{


/*45:*/


//line amail.w:659

glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n",b.name,len(b.all),b.total)
if mw!=nil{
if len(b.all)!=b.total&&b.total/100!=0&&len(b.all)%(b.total/100)!=0{
continue
}

if err:=mw.WriteAddr("0/^%s.*\\n/",escape(b.name));err!=nil{
glog.V(debug).Infof("can't write to 'addr' file: %s\n",err)
continue
}

if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if len(b.all)==b.total{
if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmt,b.name,len(b.unread),len(b.all))));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
w:=mw


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:680



/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:681

}else if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,
b.name,
len(b.unread),
len(b.all),
len(b.all)*100/b.total)));err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
}
}



/*:45*/


//line amail.w:505

}


/*44:*/


//line amail.w:604

case ev,ok:=<-ech:
glog.V(debug).Infof("an event from main window has been received: %v\n",ev)
if!ok{
ech= nil


/*9:*/


//line amail.w:118

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:609

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


//line amail.w:91

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:626

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


/*171:*/


//line amail.w:2233

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:171*/


//line amail.w:638

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


/*73:*/


//line amail.w:1005

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:73*/


//line amail.w:649

continue
}
}
mw.UnreadEvent(ev)




/*:44*/


//line amail.w:507

}
}
}()



/*:35*/


//line amail.w:69

go func(){


/*26:*/


//line amail.w:288

{
glog.V(debug).Infoln("enumerating of mailboxes")
fi,err:=rfid.Dirreadall()
if err!=nil{
glog.Errorf("can't read mailfs: %v\n",err)


/*5:*/


//line amail.w:91

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:294

return
}
for _,f:=range fi{
if f.Mode&plan9.DMDIR==plan9.DMDIR{
name:=f.Name


/*27:*/


//line amail.w:307

glog.V(debug).Infof("send a mailbox '%s' to put in the list\n",name)
bch<-name



/*:27*/


//line amail.w:300

}
}
glog.V(debug).Infoln("enumerating of mailboxes is done")
}



/*:26*/


//line amail.w:71

}()
}


/*34:*/


//line amail.w:465

glog.V(debug).Infoln("process events are specific for the list of mailboxes")
for{
select{


/*4:*/


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:469

case name:=<-bch:


/*53:*/


//line amail.w:796

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:53*/


//line amail.w:471



/*36:*/


//line amail.w:513

glog.V(debug).Infof("creating a '%s' mailbox\n",name)
box:=&mailbox{name:name,

/*22:*/


//line amail.w:241

mch:make(chan int,100),
dch:make(chan int,100),




/*:22*/



/*72:*/


//line amail.w:999

shownew:shownew,
showthreads:showthreads,
cch:make(chan bool,100),



/*:72*/



/*107:*/


//line amail.w:1461

mdch:make(chan messages,100),



/*:107*/



/*141:*/


//line amail.w:1906

rfch:make(chan*refresh,100),
irfch:make(chan*refresh,100),



/*:141*/



/*187:*/


//line amail.w:2511

ach:make(chan*struct{ids[]int;a action},100),



/*:187*/


//line amail.w:515
}
boxes= append(boxes,box)
sort.Sort(boxes)



/*:36*/


//line amail.w:472



/*70:*/


//line amail.w:979

glog.V(debug).Infoln("sending a signal to refresh all mailboxes")
rfch<-nil



/*:70*/


//line amail.w:473



/*59:*/


//line amail.w:860

go box.loop()



/*:59*/


//line amail.w:474

case d:=<-mch:
name:=d.name


/*38:*/


//line amail.w:532

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*53:*/


//line amail.w:796

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:53*/


//line amail.w:536

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:38*/


//line amail.w:477

glog.V(debug).Infof("sending '%d' to add in the '%s' mailbox\n",d.id,boxes[i].name)
boxes[i].mch<-d.id
case d:=<-dch:
name:=d.name


/*38:*/


//line amail.w:532

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*53:*/


//line amail.w:796

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:53*/


//line amail.w:536

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:38*/


//line amail.w:482

glog.V(debug).Infof("sending '%d' to delete from the '%s' mailbox\n",d.id,boxes[i].name)
boxes[i].dch<-d.id


/*7:*/


//line amail.w:104

case i:=<-wch:
wcount+= i
if wcount==0{


/*5:*/


//line amail.w:91

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:108

return
}



/*:7*/



/*87:*/


//line amail.w:1235

case d:=<-ach:
if d.m==nil{
continue
}
for name,ids:=range d.m{


/*38:*/


//line amail.w:532

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*53:*/


//line amail.w:796

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:53*/


//line amail.w:536

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:38*/


//line amail.w:1241

boxes[i].ach<-&struct{ids[]int;a action}{ids,d.a}
}



/*:87*/



/*110:*/


//line amail.w:1487

case msgs:=<-mdch:
for i,_:=range boxes{
glog.V(debug).Infof("sending %d messages to delete in the '%s' mailbox\n",len(msgs),boxes[i].name)
boxes[i].mdch<-append(messages{},msgs...)
}



/*:110*/



/*157:*/


//line amail.w:2073

case r:=<-mrfch:
for i,_:=range boxes{
glog.V(debug).Infof("sending messages to refresh in the '%s' mailbox\n",boxes[i].name)
boxes[i].rfch<-&refresh{r.flags,append(messages{},r.msgs...)}
}





/*:157*/


//line amail.w:485

}
}




/*:34*/


//line amail.w:74

}



/*:2*/



/*24:*/


//line amail.w:260

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


//line amail.w:323

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

/*211:*/


//line amail.w:2946

cids:make(map[string]*file),



/*:211*/


//line amail.w:335
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


/*94:*/


//line amail.w:1300

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



/*:94*/



/*122:*/


//line amail.w:1608

{
if _,err:=fmt.Sscanf(s,"inreplyto %s",&msg.inreplyto);err==nil{
msg.inreplyto= strings.Trim(msg.inreplyto,"<>")
continue
}
if _,err:=fmt.Sscanf(s,"messageid %s",&msg.messageid);err==nil{
msg.messageid= strings.Trim(msg.messageid,"<>")
ch:=make(chan bool)
idch<-struct{msg*message;val interface{}}{msg,ch}
if ok:=<-ch;!ok{
return nil,false,errors.New(fmt.Sprintf("a message '%s' is duplicated",msg.messageid))
}
continue
}
}



/*:122*/



/*196:*/


//line amail.w:2610

if strings.HasPrefix(s,"to "){
msg.to= split(s[len("to "):])
continue
}
if strings.HasPrefix(s,"cc "){
msg.cc= split(s[len("cc "):])
continue
}




/*:196*/


//line amail.w:349

}
msg.unread= unread
return

}



/*:30*/



/*37:*/


//line amail.w:521

func(this mailboxes)Search(name string)(int,bool){
pos:=sort.Search(len(this),
func(i int)bool{return this[i].name>=name});
if pos!=len(this)&&this[pos].name==name{
return pos,true
}
return pos,false
}



/*:37*/



/*51:*/


//line amail.w:777

func escape(s string)(res string){
for _,v:=range s{
if strings.ContainsRune("\\/[].+?()*^$",v){
res+= "\\"
}
res+= string(v)
}
return res
}



/*:51*/



/*54:*/


//line amail.w:804

func(this messages)Search(id int)(int,bool){
pos:=sort.Search(len(this),func(i int)bool{return this[i].id<=id});
if pos!=len(this)&&this[pos].id==id{
return pos,true
}
return pos,false
}



/*:54*/



/*55:*/


//line amail.w:814

func(this*messages)Insert(msg*message,pos int){
*this= append(*this,nil)
copy((*this)[pos+1:],(*this)[pos:])
(*this)[pos]= msg
}



/*:55*/



/*56:*/


//line amail.w:824

func(this*messages)SearchInsert(msg*message)(int,bool){
pos,ok:=this.Search(msg.id)
if ok{
return pos,false
}
this.Insert(msg,pos)
return pos,true
}



/*:56*/



/*57:*/


//line amail.w:837

func(this*messages)Delete(pos int)(*message,bool){
if pos<0||pos> len(*this)-1{
return nil,false
}
msg:=(*this)[pos]
*this= append((*this)[:pos],(*this)[pos+1:]...)
return msg,true
}



/*:57*/



/*58:*/


//line amail.w:850

func(this*messages)DeleteById(id int)(*message,bool){
pos,ok:=this.Search(id)
if!ok{
return nil,false
}
return this.Delete(pos)
}



/*:58*/



/*60:*/


//line amail.w:864

func(box*mailbox)loop(){
glog.V(debug).Infof("start a message loop for the '%s' mailbox\n",box.name)
counted:=false
pcount:=0
ontop:=false


/*48:*/


//line amail.w:706

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


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:725



/*61:*/


//line amail.w:889

case id:=<-box.mch:
glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n",id,box.name)
msg,new,err:=box.newMessage(id)
if err!=nil{
continue
}
if new{


/*64:*/


//line amail.w:927

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:64*/


//line amail.w:897

}
box.total++


/*65:*/


//line amail.w:935

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:65*/


//line amail.w:900



/*155:*/


//line amail.w:2061

{
glog.V(debug).Infof("print '%s/%d' at exact position\n",box.name,msg.id)
mrfch<-&refresh{seek|insert|exact,append(messages{},msg)}
}



/*:155*/


//line amail.w:901

if!box.thread{
if box.threadMode(){


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:904

var msgs messages
src:=append(messages{},root)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:907



/*151:*/


//line amail.w:2029

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)
box.rfch<-&refresh{seek|insert,msgs}
}



/*:151*/


//line amail.w:908

}else{


/*150:*/


//line amail.w:2022

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&refresh{seek|insert,append(messages{},msg)}
}



/*:150*/


//line amail.w:910

}
}


/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:913





/*:61*/



/*62:*/


//line amail.w:917

case id:=<-box.dch:
glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n",id,box.name)


/*66:*/


//line amail.w:943

if i,ok:=box.all.Search(id);ok{
msgs:=append(messages{},box.all[i])


/*67:*/


//line amail.w:952

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*136:*/


//line amail.w:1864

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{msg*message;val interface{}}{msg:msg}
}



/*:136*/


//line amail.w:961



/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:962

}
}



/*:67*/


//line amail.w:946



/*111:*/


//line amail.w:1495

mdch<-msgs



/*:111*/


//line amail.w:947

}



/*:66*/


//line amail.w:920




/*:62*/



/*74:*/


//line amail.w:1010

case<-box.cch:
glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
box.shownew= shownew
box.showthreads= showthreads
box.thread= false


/*75:*/


//line amail.w:1026

glog.V(debug).Infof("creation a window for the '%s' mailbox\n",box.name)
var err error
if box.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
if box.ech,err= box.w.EventChannel(0,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*174:*/


//line amail.w:2264

name:="Amail/"+box.name
w:=box.w


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2267




/*:174*/


//line amail.w:1037



/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1038





/*:75*/


//line amail.w:1017



/*145:*/


//line amail.w:1961

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:1964

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
var roots messages
for len(src)> 0{
msg:=src[0]


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:1972

if!roots.Check(root){


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1974

continue
}
glog.V(debug).Infof("root of thread: '%s/%d'\n",root.box.name,root.id)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:1978

}
}else{


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:1981

}
box.rfch<-&refresh{0,msgs}
}



/*:145*/


//line amail.w:1018



/*8:*/


//line amail.w:113

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:1019

}else{
glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n",box.name)
box.w.WriteCtl("dot=addr\nshow")
}



/*:74*/



/*77:*/


//line amail.w:1046

case ev,ok:=<-box.ech:
glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n",box.name,ev)
if!ok{
box.ech= nil
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":


/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:1056

box.w.Del(true)
box.w.Close()
box.w= nil


/*9:*/


//line amail.w:118

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:1060

continue
case"ShowNew":
box.thread= false
box.shownew= true
case"ShowAll":
if box.showthreads&&!counted{
continue
}
box.thread= false
box.shownew= false
case"ShowThreads":
if!counted{
continue
}
box.showthreads= true
if box.shownew==true{


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1077

continue
}
case"ShowPlain":
box.showthreads= false
if box.shownew==true{


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1083

continue
}
case"Thread":
if!counted{
continue
}
var msg*message
if len(ev.Arg)==0{


/*89:*/


//line amail.w:1250

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



/*:89*/


//line amail.w:1092

}else if num,err:=strconv.Atoi(strings.TrimSpace(ev.Arg));err!=nil{
continue
}else if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
if msg!=nil{
box.thread= true


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1100



/*177:*/


//line amail.w:2340

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:177*/


//line amail.w:1101



/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:1102



/*149:*/


//line amail.w:2014



/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:2015

var msgs messages
src:=append(messages{},root)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:2018

box.rfch<-&refresh{0,msgs}



/*:149*/


//line amail.w:1103

}
continue
case"Delmesg":


/*97:*/


//line amail.w:1334



/*96:*/


//line amail.w:1323



/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1324

if(ev.Type&goacme.Tag)==goacme.Tag&&len(ev.Arg)> 0{
s:=ev.Arg


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1327

}else if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1331




/*:96*/


//line amail.w:1335

if len(msgs)!=0{


/*100:*/


//line amail.w:1355

glog.V(debug).Infoln("sending messages to mark for deletion")
ach<-&struct{m msgmap;a action}{msgs,del}



/*:100*/


//line amail.w:1337

continue
}



/*:97*/


//line amail.w:1107

continue
case"UnDelmesg":


/*98:*/


//line amail.w:1342



/*96:*/


//line amail.w:1323



/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1324

if(ev.Type&goacme.Tag)==goacme.Tag&&len(ev.Arg)> 0{
s:=ev.Arg


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1327

}else if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1331




/*:96*/


//line amail.w:1343

if len(msgs)!=0{


/*101:*/


//line amail.w:1360

glog.V(debug).Infoln("sending messages to unmark for deletion")
ach<-&struct{m msgmap;a action}{msgs,undel}



/*:101*/


//line amail.w:1345

continue
}



/*:98*/


//line amail.w:1110

continue
case"Put":


/*105:*/


//line amail.w:1414

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


/*67:*/


//line amail.w:952

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*136:*/


//line amail.w:1864

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{msg*message;val interface{}}{msg:msg}
}



/*:136*/


//line amail.w:961



/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:962

}
}



/*:67*/


//line amail.w:1430

}
dmsgs:=msgs
for{
cmd:=fmt.Sprintf("delete %s",box.name)
c:=len(dmsgs)
if c==0{
break
}else if c> 50{
c= 50
}
for _,msg:=range dmsgs[:c]{
cmd= fmt.Sprintf("%s %d ",cmd,msg.id)
}
glog.V(debug).Infof("command to delete messages: '%s'\n",cmd)
if _,err:=f.Write([]byte(cmd));err!=nil{
glog.Errorf("can't delete messages: %v\n",err)
}
dmsgs= dmsgs[c:]
}


/*111:*/


//line amail.w:1495

mdch<-msgs



/*:111*/


//line amail.w:1450


f.Close()





/*:105*/


//line amail.w:1113

continue
case"Mail":
var msg*message


/*232:*/


//line amail.w:3348

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*233:*/


//line amail.w:3360

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


/*237:*/


//line amail.w:3451

{


/*240:*/


//line amail.w:3579



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3580




/*:240*/


//line amail.w:3453

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
case"to":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3478

to= append(to,f...)
case"cc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3481

cc= append(cc,f...)
case"bcc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3484

bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%s",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3499

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



/*:237*/


//line amail.w:3375

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:233*/


//line amail.w:3357




/*:232*/


//line amail.w:1117

name:=fmt.Sprintf("Amail/%s/New",box.name)


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:1119



/*242:*/


//line amail.w:3590

writeSignature(w,box)



/*:242*/


//line amail.w:1120

continue
case"Seen":


/*112:*/


//line amail.w:1499



/*96:*/


//line amail.w:1323



/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1324

if(ev.Type&goacme.Tag)==goacme.Tag&&len(ev.Arg)> 0{
s:=ev.Arg


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1327

}else if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1331




/*:96*/


//line amail.w:1500

if len(msgs)!=0{


/*114:*/


//line amail.w:1511

glog.V(debug).Infoln("sending messages to mark them seen")
ach<-&struct{m msgmap;a action}{msgs,seen}



/*:114*/


//line amail.w:1502

continue
}



/*:112*/


//line amail.w:1123

continue
case"Search":
glog.V(debug).Infof("search argument: '%s'\n",ev.Arg)


/*184:*/


//line amail.w:2455

{
msgs:=box.search(ev.Arg)


/*177:*/


//line amail.w:2340

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:177*/


//line amail.w:2458



/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:2459

name:=fmt.Sprintf("Amail/%s/Search(%s)",box.name,strings.Replace(ev.Arg," ","␣",-1))
w:=box.w
box.thread= false
box.shownew= true
box.showthreads= false


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2465

glog.V(debug).Infof("len of msgs: %v\n",len(msgs))
box.rfch<-&refresh{0,msgs}
}



/*:184*/


//line amail.w:1127

continue
default:
box.w.UnreadEvent(ev)
continue
}


/*174:*/


//line amail.w:2264

name:="Amail/"+box.name
w:=box.w


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2267




/*:174*/


//line amail.w:1133



/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1134



/*177:*/


//line amail.w:2340

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:177*/


//line amail.w:1135



/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:1136



/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:1137



/*145:*/


//line amail.w:1961

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:1964

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
var roots messages
for len(src)> 0{
msg:=src[0]


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:1972

if!roots.Check(root){


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1974

continue
}
glog.V(debug).Infof("root of thread: '%s/%d'\n",root.box.name,root.id)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:1978

}
}else{


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:1981

}
box.rfch<-&refresh{0,msgs}
}



/*:145*/


//line amail.w:1138

continue
}else if(ev.Type&goacme.Look)==goacme.Look{


/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1141

if(ev.Type&goacme.Tag)==goacme.Tag{
s:=ev.Text


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1144

}else{


/*78:*/


//line amail.w:1159

glog.V(debug).Infof("event: %v\n",ev)
if err:=box.w.WriteAddr("#%d,#%d",ev.Begin,ev.End);err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1164




/*:78*/


//line amail.w:1146

}
if len(msgs)!=0{


/*88:*/


//line amail.w:1246

ach<-&struct{m msgmap;a action}{msgs,view}



/*:88*/


//line amail.w:1149

continue
}
}
box.w.UnreadEvent(ev)



/*:77*/



/*108:*/


//line amail.w:1466

case m:=<-box.mdch:
if box.w==nil{
continue
}
glog.V(debug).Infof("%d messages were received to be deleted from the '%s' mailbox\n",len(m),box.name)
for _,msg:=range m{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:1473

if box.threadMode(){


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:1475



/*170:*/


//line amail.w:2215

{
if len(children)!=0{
var msgs messages
var src messages
for _,msg:=range children{
if msg.box!=box{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:2222

}else{
root:=msg


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:2225

}
}


/*153:*/


//line amail.w:2043

{
if len(msgs)!=0{
glog.V(debug).Infof("inform the '%s' mailbox to refresh messages\n",box.name)
box.rfch<-&refresh{seek,msgs}
}
}



/*:153*/


//line amail.w:2228

}
}



/*:170*/


//line amail.w:1476

}
}


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:1479




/*:108*/



/*142:*/


//line amail.w:1916

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


/*164:*/


//line amail.w:2125

{
glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with flags: %v\n",box.name,len(v.msgs),v.flags)
f,err:=box.w.File("data")
if err!=nil{
glog.Errorf("can't open 'data' file of the '%s' messagebox: %v\n",box.name,err)
continue
}
if(v.flags&seek)==seek{


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2134

msg:=v.msgs[0]


/*180:*/


//line amail.w:2364



/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:2365



/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2366

glog.V(debug).Infof("composed message addr '%s' in the '%s' mailbox\n",addr,box.name)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("the '%d' message is not found in the window of the '%s' mailbox\n",msg.id,box.name)
if(v.flags&insert)==0{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2371

}
if box.threadMode(){


/*183:*/


//line amail.w:2417



/*131:*/


//line amail.w:1799

var parent*message
{
ch:=make(chan parentmsg)
glog.V(debug).Infof("getting parent for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
parent= <-ch
}



/*:131*/


//line amail.w:2418

if parent!=nil{
glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n",msg.id)
m:=msg
msg= parent
found:=false
for!found{


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2425

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


/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2440

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
if(v.flags&exact)==exact{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2445

}
}
}else if(v.flags&exact)==exact{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2449

}else if err:=box.w.WriteAddr(bof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %v\n",bof,box.name,err)
}



/*:183*/


//line amail.w:2374

}else if msg.box!=box{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2376

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


/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2389

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
}
}
}



/*:180*/


//line amail.w:2136

}else if err:=box.w.WriteAddr(eof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",eof,err)
continue
}
w:=box.w
glog.V(debug).Infof("printing of messages of the '%s' mailbox\n",box.name)
buf:=make([]byte,0,0x8000)


/*167:*/


//line amail.w:2176

c:=0
for len(v.msgs)> 0&&c<100{
msg:=v.msgs[0]
glog.V(debug).Infof("printing of '%s/%d' message with in the '%s' mailbox\n",msg.box.name,msg.id,box.name)
if box.threadMode(){


/*169:*/


//line amail.w:2204

{


/*135:*/


//line amail.w:1854

var level int
{
ch:=make(chan int)
glog.V(debug).Infof("getting root for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
level= <-ch
}



/*:135*/


//line amail.w:2206

for;level> 0;level--{
buf= append(buf,levelmark...)
}
}



/*:169*/


//line amail.w:2182

}
c++


/*91:*/


//line amail.w:1277

glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n",msg.id,box.name)
buf= append(buf,fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n",
func()string{if msg.deleted{return deleted};return""}(),
func()string{if msg.box!=box{return fmt.Sprintf("%s/",msg.box.name)};return""}(),
msg.id,
func()string{if msg.unread{return newmark};return""}(),
msg.from,
msg.date,
msg.subject)...)




/*:91*/


//line amail.w:2185

v.msgs= v.msgs[1:]
if(v.flags&seek)==seek{
break
}
}
pcount+= c



/*:167*/


//line amail.w:2144

if _,err:=f.Write(buf);err!=nil{
glog.Errorf("can't write to 'data' file of the '%s' messagebox: %v\n",box.name,err)
}


/*166:*/


//line amail.w:2162

if!ontop{
glog.V(debug).Infof("pcount: %v, ontop: %v\n",pcount,ontop)


/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:2165

if pcount>=100{
ontop= true
}
}




/*:166*/


//line amail.w:2148



/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2149

}



/*:164*/


//line amail.w:1930





/*:142*/



/*189:*/


//line amail.w:2521

case d:=<-box.ach:
switch d.a{
case view:
var msgs messages
for _,id:=range d.ids{
glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n",id,box.name)
p,ok:=box.all.Search(id)
if!ok{
glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n",id,box.name)
continue
}
msg:=box.all[p]
if msg.w==nil{
if msg.unread{


/*190:*/


//line amail.w:2553

msg.unread= false
box.unread.DeleteById(id)



/*:190*/


//line amail.w:2536



/*191:*/


//line amail.w:2559

if!box.thread&&box.shownew{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:2561



/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2562

}
msgs= append(msgs,msg)





/*:191*/


//line amail.w:2537



/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:2538

}
if err:=msg.open();err!=nil{
continue
}
}else{
glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n",id,box.name)
msg.w.WriteCtl("dot=addr\nshow")
}
}


/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:2548



/*102:*/


//line amail.w:1365

case del:
var msgs messages
for _,id:=range d.ids{


/*103:*/


//line amail.w:1380

if p,ok:=box.all.Search(id);ok{
if box.all[p].deleted{
continue
}
box.all[p].deleted= true
box.deleted++
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*199:*/


//line amail.w:2675

msg.writeTag()



/*:199*/


//line amail.w:1390

}
glog.V(debug).Infof("the '%v' message is marked for deletion\n",id)
}



/*:103*/


//line amail.w:1369

}


/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:1371

case undel:
var msgs messages
for _,id:=range d.ids{


/*104:*/


//line amail.w:1396

if p,ok:=box.all.Search(id);ok{
if!box.all[p].deleted{
continue
}
box.all[p].deleted= false
box.deleted--
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*199:*/


//line amail.w:2675

msg.writeTag()



/*:199*/


//line amail.w:1406

}
glog.V(debug).Infof("the '%v' message is unmarked for deletion\n",id)

}



/*:104*/


//line amail.w:1375

}


/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:1377




/*:102*/



/*115:*/


//line amail.w:1517

case seen:
f,err:=box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl' file of the '%s' messagebox: %v\n",box.name,err)
continue
}
for{
c:=len(d.ids)
if c==0{
break
}else if c> 50{
c= 50
}
ids:=d.ids[0:c]
d.ids= d.ids[c:]
var ms messages
for _,id:=range ids{
p,ok:=box.all.Search(id)
if!ok||!box.all[p].unread{
continue
}
ms= append(ms,box.all[p])
}
cmd:="read"
for _,v:=range ms{
cmd+= fmt.Sprintf(" %d",v.id)
}
if _,err:=f.Write([]byte(cmd));err!=nil{
glog.Errorf("can't write to 'ctl' file of the '%s' messagebox: %v\n",box.name,err)
}
var msgs messages
for _,msg:=range ms{
id:=msg.id


/*190:*/


//line amail.w:2553

msg.unread= false
box.unread.DeleteById(id)



/*:190*/


//line amail.w:1554



/*191:*/


//line amail.w:2559

if!box.thread&&box.shownew{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:2561



/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2562

}
msgs= append(msgs,msg)





/*:191*/


//line amail.w:1555

}


/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:1557



/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:1558

}
f.Close()



/*:115*/


//line amail.w:2549

}



/*:189*/


//line amail.w:726

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


//line amail.w:307

glog.V(debug).Infof("send a mailbox '%s' to put in the list\n",name)
bch<-name



/*:27*/


//line amail.w:738

continue
}
if msg,new,err:=box.newMessage(id);err==nil{
if new{


/*64:*/


//line amail.w:927

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:64*/


//line amail.w:743

}


/*65:*/


//line amail.w:935

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:65*/


//line amail.w:745

}else{
glog.V(debug).Infof("can't create a new '%d' message in the '%s' mailbox: %v\n",id,box.name,err)
box.total--
continue
}


/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:751



/*161:*/


//line amail.w:2097

if!box.threadMode(){


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:2099

if len(src)!=0&&box.pos<len(src)&&len(src)%500==0{
glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n",box.name,len(src)-box.pos)
msgs:=append(messages{},src[box.pos:len(src)]...)


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:2103

box.rfch<-&refresh{0,msgs}
}
}



/*:161*/


//line amail.w:752

}
}


/*162:*/


//line amail.w:2109

if!box.threadMode(){


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:2111

if box.pos<len(src){
glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n",box.name,len(src)-box.pos)
msgs:=append(messages{},src[box.pos:len(src)]...)


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:2115

box.rfch<-&refresh{0,msgs}
}
}



/*:162*/


//line amail.w:755

}



/*:48*/


//line amail.w:870

counted= true


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:872

if box.threadMode(){


/*145:*/


//line amail.w:1961

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:1964

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
var roots messages
for len(src)> 0{
msg:=src[0]


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:1972

if!roots.Check(root){


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1974

continue
}
glog.V(debug).Infof("root of thread: '%s/%d'\n",root.box.name,root.id)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:1978

}
}else{


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:1981

}
box.rfch<-&refresh{0,msgs}
}



/*:145*/


//line amail.w:874

}
defer glog.V(debug).Infof("a message loop of the '%s' mailbox is done\n",box.name)
for{
select{


/*4:*/


//line amail.w:85

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:879



/*61:*/


//line amail.w:889

case id:=<-box.mch:
glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n",id,box.name)
msg,new,err:=box.newMessage(id)
if err!=nil{
continue
}
if new{


/*64:*/


//line amail.w:927

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:64*/


//line amail.w:897

}
box.total++


/*65:*/


//line amail.w:935

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:65*/


//line amail.w:900



/*155:*/


//line amail.w:2061

{
glog.V(debug).Infof("print '%s/%d' at exact position\n",box.name,msg.id)
mrfch<-&refresh{seek|insert|exact,append(messages{},msg)}
}



/*:155*/


//line amail.w:901

if!box.thread{
if box.threadMode(){


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:904

var msgs messages
src:=append(messages{},root)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:907



/*151:*/


//line amail.w:2029

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)
box.rfch<-&refresh{seek|insert,msgs}
}



/*:151*/


//line amail.w:908

}else{


/*150:*/


//line amail.w:2022

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&refresh{seek|insert,append(messages{},msg)}
}



/*:150*/


//line amail.w:910

}
}


/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:913





/*:61*/



/*62:*/


//line amail.w:917

case id:=<-box.dch:
glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n",id,box.name)


/*66:*/


//line amail.w:943

if i,ok:=box.all.Search(id);ok{
msgs:=append(messages{},box.all[i])


/*67:*/


//line amail.w:952

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*136:*/


//line amail.w:1864

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{msg*message;val interface{}}{msg:msg}
}



/*:136*/


//line amail.w:961



/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:962

}
}



/*:67*/


//line amail.w:946



/*111:*/


//line amail.w:1495

mdch<-msgs



/*:111*/


//line amail.w:947

}



/*:66*/


//line amail.w:920




/*:62*/



/*74:*/


//line amail.w:1010

case<-box.cch:
glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
box.shownew= shownew
box.showthreads= showthreads
box.thread= false


/*75:*/


//line amail.w:1026

glog.V(debug).Infof("creation a window for the '%s' mailbox\n",box.name)
var err error
if box.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
if box.ech,err= box.w.EventChannel(0,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*174:*/


//line amail.w:2264

name:="Amail/"+box.name
w:=box.w


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2267




/*:174*/


//line amail.w:1037



/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1038





/*:75*/


//line amail.w:1017



/*145:*/


//line amail.w:1961

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:1964

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
var roots messages
for len(src)> 0{
msg:=src[0]


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:1972

if!roots.Check(root){


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1974

continue
}
glog.V(debug).Infof("root of thread: '%s/%d'\n",root.box.name,root.id)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:1978

}
}else{


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:1981

}
box.rfch<-&refresh{0,msgs}
}



/*:145*/


//line amail.w:1018



/*8:*/


//line amail.w:113

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:1019

}else{
glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n",box.name)
box.w.WriteCtl("dot=addr\nshow")
}



/*:74*/



/*77:*/


//line amail.w:1046

case ev,ok:=<-box.ech:
glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n",box.name,ev)
if!ok{
box.ech= nil
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":


/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:1056

box.w.Del(true)
box.w.Close()
box.w= nil


/*9:*/


//line amail.w:118

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:1060

continue
case"ShowNew":
box.thread= false
box.shownew= true
case"ShowAll":
if box.showthreads&&!counted{
continue
}
box.thread= false
box.shownew= false
case"ShowThreads":
if!counted{
continue
}
box.showthreads= true
if box.shownew==true{


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1077

continue
}
case"ShowPlain":
box.showthreads= false
if box.shownew==true{


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1083

continue
}
case"Thread":
if!counted{
continue
}
var msg*message
if len(ev.Arg)==0{


/*89:*/


//line amail.w:1250

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



/*:89*/


//line amail.w:1092

}else if num,err:=strconv.Atoi(strings.TrimSpace(ev.Arg));err!=nil{
continue
}else if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
if msg!=nil{
box.thread= true


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1100



/*177:*/


//line amail.w:2340

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:177*/


//line amail.w:1101



/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:1102



/*149:*/


//line amail.w:2014



/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:2015

var msgs messages
src:=append(messages{},root)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:2018

box.rfch<-&refresh{0,msgs}



/*:149*/


//line amail.w:1103

}
continue
case"Delmesg":


/*97:*/


//line amail.w:1334



/*96:*/


//line amail.w:1323



/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1324

if(ev.Type&goacme.Tag)==goacme.Tag&&len(ev.Arg)> 0{
s:=ev.Arg


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1327

}else if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1331




/*:96*/


//line amail.w:1335

if len(msgs)!=0{


/*100:*/


//line amail.w:1355

glog.V(debug).Infoln("sending messages to mark for deletion")
ach<-&struct{m msgmap;a action}{msgs,del}



/*:100*/


//line amail.w:1337

continue
}



/*:97*/


//line amail.w:1107

continue
case"UnDelmesg":


/*98:*/


//line amail.w:1342



/*96:*/


//line amail.w:1323



/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1324

if(ev.Type&goacme.Tag)==goacme.Tag&&len(ev.Arg)> 0{
s:=ev.Arg


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1327

}else if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1331




/*:96*/


//line amail.w:1343

if len(msgs)!=0{


/*101:*/


//line amail.w:1360

glog.V(debug).Infoln("sending messages to unmark for deletion")
ach<-&struct{m msgmap;a action}{msgs,undel}



/*:101*/


//line amail.w:1345

continue
}



/*:98*/


//line amail.w:1110

continue
case"Put":


/*105:*/


//line amail.w:1414

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


/*67:*/


//line amail.w:952

{
if msg,ok:=box.all.Delete(i);ok{
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",msg.id,box.name)
box.unread.DeleteById(msg.id)
box.total--
if msg.deleted{
box.deleted--
}


/*136:*/


//line amail.w:1864

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{msg*message;val interface{}}{msg:msg}
}



/*:136*/


//line amail.w:961



/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:962

}
}



/*:67*/


//line amail.w:1430

}
dmsgs:=msgs
for{
cmd:=fmt.Sprintf("delete %s",box.name)
c:=len(dmsgs)
if c==0{
break
}else if c> 50{
c= 50
}
for _,msg:=range dmsgs[:c]{
cmd= fmt.Sprintf("%s %d ",cmd,msg.id)
}
glog.V(debug).Infof("command to delete messages: '%s'\n",cmd)
if _,err:=f.Write([]byte(cmd));err!=nil{
glog.Errorf("can't delete messages: %v\n",err)
}
dmsgs= dmsgs[c:]
}


/*111:*/


//line amail.w:1495

mdch<-msgs



/*:111*/


//line amail.w:1450


f.Close()





/*:105*/


//line amail.w:1113

continue
case"Mail":
var msg*message


/*232:*/


//line amail.w:3348

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*233:*/


//line amail.w:3360

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


/*237:*/


//line amail.w:3451

{


/*240:*/


//line amail.w:3579



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3580




/*:240*/


//line amail.w:3453

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
case"to":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3478

to= append(to,f...)
case"cc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3481

cc= append(cc,f...)
case"bcc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3484

bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%s",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3499

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



/*:237*/


//line amail.w:3375

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:233*/


//line amail.w:3357




/*:232*/


//line amail.w:1117

name:=fmt.Sprintf("Amail/%s/New",box.name)


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:1119



/*242:*/


//line amail.w:3590

writeSignature(w,box)



/*:242*/


//line amail.w:1120

continue
case"Seen":


/*112:*/


//line amail.w:1499



/*96:*/


//line amail.w:1323



/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1324

if(ev.Type&goacme.Tag)==goacme.Tag&&len(ev.Arg)> 0{
s:=ev.Arg


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1327

}else if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1331




/*:96*/


//line amail.w:1500

if len(msgs)!=0{


/*114:*/


//line amail.w:1511

glog.V(debug).Infoln("sending messages to mark them seen")
ach<-&struct{m msgmap;a action}{msgs,seen}



/*:114*/


//line amail.w:1502

continue
}



/*:112*/


//line amail.w:1123

continue
case"Search":
glog.V(debug).Infof("search argument: '%s'\n",ev.Arg)


/*184:*/


//line amail.w:2455

{
msgs:=box.search(ev.Arg)


/*177:*/


//line amail.w:2340

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:177*/


//line amail.w:2458



/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:2459

name:=fmt.Sprintf("Amail/%s/Search(%s)",box.name,strings.Replace(ev.Arg," ","␣",-1))
w:=box.w
box.thread= false
box.shownew= true
box.showthreads= false


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2465

glog.V(debug).Infof("len of msgs: %v\n",len(msgs))
box.rfch<-&refresh{0,msgs}
}



/*:184*/


//line amail.w:1127

continue
default:
box.w.UnreadEvent(ev)
continue
}


/*174:*/


//line amail.w:2264

name:="Amail/"+box.name
w:=box.w


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2267




/*:174*/


//line amail.w:1133



/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:1134



/*177:*/


//line amail.w:2340

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:177*/


//line amail.w:1135



/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:1136



/*159:*/


//line amail.w:2087

box.pos= 0
ontop= false



/*:159*/



/*168:*/


//line amail.w:2194

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*refresh,100)
pcount= 0
ontop= false
}



/*:168*/


//line amail.w:1137



/*145:*/


//line amail.w:1961

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:1964

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
var roots messages
for len(src)> 0{
msg:=src[0]


/*133:*/


//line amail.w:1827

root:=msg
{
ch:=make(chan rootmsg)
glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n",msg.messageid,msg.box.name,msg.id)
idch<-struct{msg*message;val interface{}}{msg,ch}
root= <-ch
}



/*:133*/


//line amail.w:1972

if!roots.Check(root){


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1974

continue
}
glog.V(debug).Infof("root of thread: '%s/%d'\n",root.box.name,root.id)


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:1978

}
}else{


/*160:*/


//line amail.w:2092

box.pos= len(src)



/*:160*/


//line amail.w:1981

}
box.rfch<-&refresh{0,msgs}
}



/*:145*/


//line amail.w:1138

continue
}else if(ev.Type&goacme.Look)==goacme.Look{


/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:1141

if(ev.Type&goacme.Tag)==goacme.Tag{
s:=ev.Text


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1144

}else{


/*78:*/


//line amail.w:1159

glog.V(debug).Infof("event: %v\n",ev)
if err:=box.w.WriteAddr("#%d,#%d",ev.Begin,ev.End);err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else


/*79:*/


//line amail.w:1167

if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*80:*/


//line amail.w:1183

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


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:1198

break
}
}
}



/*:80*/


//line amail.w:1173

if err==io.EOF{
break
}
}
}




/*:79*/


//line amail.w:1164




/*:78*/


//line amail.w:1146

}
if len(msgs)!=0{


/*88:*/


//line amail.w:1246

ach<-&struct{m msgmap;a action}{msgs,view}



/*:88*/


//line amail.w:1149

continue
}
}
box.w.UnreadEvent(ev)



/*:77*/



/*108:*/


//line amail.w:1466

case m:=<-box.mdch:
if box.w==nil{
continue
}
glog.V(debug).Infof("%d messages were received to be deleted from the '%s' mailbox\n",len(m),box.name)
for _,msg:=range m{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:1473

if box.threadMode(){


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:1475



/*170:*/


//line amail.w:2215

{
if len(children)!=0{
var msgs messages
var src messages
for _,msg:=range children{
if msg.box!=box{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:2222

}else{
root:=msg


/*146:*/


//line amail.w:1988

msgs= append(msgs,root)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:1990

msgs,src= getchildren(root,msgs,src)



/*:146*/


//line amail.w:2225

}
}


/*153:*/


//line amail.w:2043

{
if len(msgs)!=0{
glog.V(debug).Infof("inform the '%s' mailbox to refresh messages\n",box.name)
box.rfch<-&refresh{seek,msgs}
}
}



/*:153*/


//line amail.w:2228

}
}



/*:170*/


//line amail.w:1476

}
}


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:1479




/*:108*/



/*142:*/


//line amail.w:1916

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


/*164:*/


//line amail.w:2125

{
glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with flags: %v\n",box.name,len(v.msgs),v.flags)
f,err:=box.w.File("data")
if err!=nil{
glog.Errorf("can't open 'data' file of the '%s' messagebox: %v\n",box.name,err)
continue
}
if(v.flags&seek)==seek{


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2134

msg:=v.msgs[0]


/*180:*/


//line amail.w:2364



/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:2365



/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2366

glog.V(debug).Infof("composed message addr '%s' in the '%s' mailbox\n",addr,box.name)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("the '%d' message is not found in the window of the '%s' mailbox\n",msg.id,box.name)
if(v.flags&insert)==0{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2371

}
if box.threadMode(){


/*183:*/


//line amail.w:2417



/*131:*/


//line amail.w:1799

var parent*message
{
ch:=make(chan parentmsg)
glog.V(debug).Infof("getting parent for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
parent= <-ch
}



/*:131*/


//line amail.w:2418

if parent!=nil{
glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n",msg.id)
m:=msg
msg= parent
found:=false
for!found{


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2425

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


/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2440

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
if(v.flags&exact)==exact{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2445

}
}
}else if(v.flags&exact)==exact{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2449

}else if err:=box.w.WriteAddr(bof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %v\n",bof,box.name,err)
}



/*:183*/


//line amail.w:2374

}else if msg.box!=box{


/*181:*/


//line amail.w:2398

glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n",v.msgs[0].id,box.name)
v.msgs= v.msgs[1:]


/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2401

continue




/*:181*/


//line amail.w:2376

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


/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2389

addr+= eol
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
}
}
}



/*:180*/


//line amail.w:2136

}else if err:=box.w.WriteAddr(eof);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",eof,err)
continue
}
w:=box.w
glog.V(debug).Infof("printing of messages of the '%s' mailbox\n",box.name)
buf:=make([]byte,0,0x8000)


/*167:*/


//line amail.w:2176

c:=0
for len(v.msgs)> 0&&c<100{
msg:=v.msgs[0]
glog.V(debug).Infof("printing of '%s/%d' message with in the '%s' mailbox\n",msg.box.name,msg.id,box.name)
if box.threadMode(){


/*169:*/


//line amail.w:2204

{


/*135:*/


//line amail.w:1854

var level int
{
ch:=make(chan int)
glog.V(debug).Infof("getting root for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
level= <-ch
}



/*:135*/


//line amail.w:2206

for;level> 0;level--{
buf= append(buf,levelmark...)
}
}



/*:169*/


//line amail.w:2182

}
c++


/*91:*/


//line amail.w:1277

glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n",msg.id,box.name)
buf= append(buf,fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n",
func()string{if msg.deleted{return deleted};return""}(),
func()string{if msg.box!=box{return fmt.Sprintf("%s/",msg.box.name)};return""}(),
msg.id,
func()string{if msg.unread{return newmark};return""}(),
msg.from,
msg.date,
msg.subject)...)




/*:91*/


//line amail.w:2185

v.msgs= v.msgs[1:]
if(v.flags&seek)==seek{
break
}
}
pcount+= c



/*:167*/


//line amail.w:2144

if _,err:=f.Write(buf);err!=nil{
glog.Errorf("can't write to 'data' file of the '%s' messagebox: %v\n",box.name,err)
}


/*166:*/


//line amail.w:2162

if!ontop{
glog.V(debug).Infof("pcount: %v, ontop: %v\n",pcount,ontop)


/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:2165

if pcount>=100{
ontop= true
}
}




/*:166*/


//line amail.w:2148



/*165:*/


//line amail.w:2153

if len(v.msgs)> 0{
box.rfch<-&refresh{v.flags,v.msgs}
}else{


/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2157

}



/*:165*/


//line amail.w:2149

}



/*:164*/


//line amail.w:1930





/*:142*/



/*189:*/


//line amail.w:2521

case d:=<-box.ach:
switch d.a{
case view:
var msgs messages
for _,id:=range d.ids{
glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n",id,box.name)
p,ok:=box.all.Search(id)
if!ok{
glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n",id,box.name)
continue
}
msg:=box.all[p]
if msg.w==nil{
if msg.unread{


/*190:*/


//line amail.w:2553

msg.unread= false
box.unread.DeleteById(id)



/*:190*/


//line amail.w:2536



/*191:*/


//line amail.w:2559

if!box.thread&&box.shownew{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:2561



/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2562

}
msgs= append(msgs,msg)





/*:191*/


//line amail.w:2537



/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:2538

}
if err:=msg.open();err!=nil{
continue
}
}else{
glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n",id,box.name)
msg.w.WriteCtl("dot=addr\nshow")
}
}


/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:2548



/*102:*/


//line amail.w:1365

case del:
var msgs messages
for _,id:=range d.ids{


/*103:*/


//line amail.w:1380

if p,ok:=box.all.Search(id);ok{
if box.all[p].deleted{
continue
}
box.all[p].deleted= true
box.deleted++
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*199:*/


//line amail.w:2675

msg.writeTag()



/*:199*/


//line amail.w:1390

}
glog.V(debug).Infof("the '%v' message is marked for deletion\n",id)
}



/*:103*/


//line amail.w:1369

}


/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:1371

case undel:
var msgs messages
for _,id:=range d.ids{


/*104:*/


//line amail.w:1396

if p,ok:=box.all.Search(id);ok{
if!box.all[p].deleted{
continue
}
box.all[p].deleted= false
box.deleted--
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
msg:=box.all[p]


/*199:*/


//line amail.w:2675

msg.writeTag()



/*:199*/


//line amail.w:1406

}
glog.V(debug).Infof("the '%v' message is unmarked for deletion\n",id)

}



/*:104*/


//line amail.w:1375

}


/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:1377




/*:102*/



/*115:*/


//line amail.w:1517

case seen:
f,err:=box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl' file of the '%s' messagebox: %v\n",box.name,err)
continue
}
for{
c:=len(d.ids)
if c==0{
break
}else if c> 50{
c= 50
}
ids:=d.ids[0:c]
d.ids= d.ids[c:]
var ms messages
for _,id:=range ids{
p,ok:=box.all.Search(id)
if!ok||!box.all[p].unread{
continue
}
ms= append(ms,box.all[p])
}
cmd:="read"
for _,v:=range ms{
cmd+= fmt.Sprintf(" %d",v.id)
}
if _,err:=f.Write([]byte(cmd));err!=nil{
glog.Errorf("can't write to 'ctl' file of the '%s' messagebox: %v\n",box.name,err)
}
var msgs messages
for _,msg:=range ms{
id:=msg.id


/*190:*/


//line amail.w:2553

msg.unread= false
box.unread.DeleteById(id)



/*:190*/


//line amail.w:1554



/*191:*/


//line amail.w:2559

if!box.thread&&box.shownew{


/*193:*/


//line amail.w:2582

box.eraseMessage(msg)




/*:193*/


//line amail.w:2561



/*192:*/


//line amail.w:2569

{
glog.V(debug).Infof("box.deleted:%d\n",box.deleted)


/*175:*/


//line amail.w:2270

box.writeTag(counted)



/*:175*/


//line amail.w:2572

w:=box.w
if box.deleted==0{


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2575

}else{


/*50:*/


//line amail.w:768

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2577

}
}



/*:192*/


//line amail.w:2562

}
msgs= append(msgs,msg)





/*:191*/


//line amail.w:1555

}


/*69:*/


//line amail.w:973

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:69*/


//line amail.w:1557



/*154:*/


//line amail.w:2052

{
if len(msgs)!=0{
glog.V(debug).Infoln("refresh messages\n")
mrfch<-&refresh{seek,msgs}
}
}



/*:154*/


//line amail.w:1558

}
f.Close()



/*:115*/


//line amail.w:2549

}



/*:189*/


//line amail.w:880

}
}
}



/*:60*/



/*68:*/


//line amail.w:967

func(box*mailbox)threadMode()bool{
return box.thread||box.showthreads&&!box.shownew
}



/*:68*/



/*126:*/


//line amail.w:1720

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




/*:126*/



/*144:*/


//line amail.w:1944

func(this*messages)Check(msg*message)bool{
pos:=sort.Search(len(*this),func(i int)bool{return(*this)[i].messageid<=msg.messageid});
if pos!=len(*this)&&(*this)[pos].messageid==msg.messageid{
return false
}
*this= append(*this,nil)
copy((*this)[pos+1:],(*this)[pos:])
(*this)[pos]= msg
return true
}




/*:144*/



/*148:*/


//line amail.w:2002

func getchildren(msg*message,dst messages,src messages)(messages,messages){


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2004

for _,msg:=range children{
dst= append(dst,msg)


/*147:*/


//line amail.w:1994

if p,ok:=src.Search(msg.id);ok{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}



/*:147*/


//line amail.w:2007

dst,src= getchildren(msg,dst,src)
}
return dst,src
}



/*:148*/



/*172:*/


//line amail.w:2243

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



/*:172*/



/*176:*/


//line amail.w:2274

func(box*mailbox)writeTag(counted bool){
glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)


/*143:*/


//line amail.w:1934

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:143*/


//line amail.w:2277

if err:=writeTag(box.w,fmt.Sprintf(" %sMail %s%s%s%s%s%sSearch ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if len(src)> 0{
return"Delmesg "
}
return""
}(),
func()string{
if len(src)> 0&&box.deleted> 0{
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
}else if len(src)> 0&&counted&&(box.shownew||!box.showthreads){
return"Thread "
}
return""
}(),
func()string{
if box.showthreads&&!counted{
return""
}
if box.shownew{
return"ShowAll "
}else{
return"ShowNew "
}
}(),
func()string{
if box.showthreads{
return"ShowPlain "
}else if counted{
return"ShowThreads "
}else{
return""
}
}(),
func()string{
if len(src)> 0{
return"Seen "
}
return""
}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}
}



/*:176*/



/*178:*/


//line amail.w:2345

func clean(w*goacme.Window){
if err:=w.WriteAddr(wholefile);err!=nil{
glog.Errorf("can't write '%s' to 'addr' file: %s\n",wholefile,err)
}else if data,err:=w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if _,err:=data.Write([]byte(""));err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
}
}



/*:178*/



/*185:*/


//line amail.w:2471

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



/*:185*/



/*194:*/


//line amail.w:2587

func(box*mailbox)eraseMessage(msg*message){
if box.w==nil{
return
}
glog.V(debug).Infof("removing the '%d' message of the '%s' mailbox from the '%s' mailbox\n",
msg.id,msg.box.name,box.name)


/*182:*/


//line amail.w:2406

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box!=msg.box{return escape(msg.box.name+"/")};return""}(),
msg.id,
escape(newmark))



/*:182*/


//line amail.w:2594

if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n",addr,box.name,err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file of the box '%s': %s\n",box.name,err)
}else if _,err:=data.Write([]byte{});err!=nil{
glog.Errorf("can't write to 'data' file of the box '%s': %s\n",box.name,err)
}
}



/*:194*/



/*197:*/


//line amail.w:2623

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



/*:197*/



/*198:*/


//line amail.w:2640

func(msg*message)open()(err error){
glog.V(debug).Infof("open: trying to open '%d' directory\n",msg.id)
bfid,err:=msg.box.fid.Walk(fmt.Sprintf("%d",msg.id))
if err!=nil{
glog.Errorf("can't walk to '%s/%d': %v\n",msg.box.name,msg.id,err)
return err
}
defer bfid.Close()
isnew:=msg.w==nil
if isnew{
if msg.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
return err
}
}else{


/*207:*/


//line amail.w:2911

glog.V(debug).Infof("clean the '%s/%d' message's window\n",msg.box.name,msg.id)
clean(msg.w)



/*:207*/


//line amail.w:2656

}
buf:=make([]byte,0,0x8000)


/*205:*/


//line amail.w:2778

{
glog.V(debug).Infof("composing a header of the '%d' message\n",msg.id)
buf= append(buf,fmt.Sprintf("From: %s\nDate: %s\nTo: %s\n%sSubject: %s\n\n\n",
msg.from,msg.date,strings.Join(msg.to,", "),
func()string{if len(msg.cc)!=0{return fmt.Sprintf("CC: %s\n",strings.Join(msg.cc,", "))};return""}(),
msg.subject)...)
}



/*:205*/


//line amail.w:2659



/*212:*/


//line amail.w:2952

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


/*216:*/


//line amail.w:3088



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3089




/*:216*/


//line amail.w:2994

for _,v:=range msg.files{
buf= append(buf,fmt.Sprintf("\n===> %s (%s)\n",v.path,v.mimetype)...)
buf= append(buf,fmt.Sprintf("\t9p read %s/%s/%d/%sbody > '%s/%s'\n",srv,msg.box.name,msg.id,v.path,home,v.name)...)
}
}



/*:212*/


//line amail.w:2660

w:=msg.w
name:=fmt.Sprintf("Amail/%s/%d",msg.box.name,msg.id)


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:2663



/*199:*/


//line amail.w:2675

msg.writeTag()



/*:199*/


//line amail.w:2664

w.Write(buf)


/*49:*/


//line amail.w:759

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:49*/


//line amail.w:2666



/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:2667

if isnew{


/*206:*/


//line amail.w:2788

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


/*152:*/


//line amail.w:2036

{
glog.V(debug).Infof("refresh a message '%d'\n",msg.id)
mrfch<-&refresh{seek,append(messages{},msg)}
}



/*:152*/


//line amail.w:2812

return
}
continue
case"UnDelmesg":
if msg.deleted{
msg.deleted= false
msg.box.deleted--


/*199:*/


//line amail.w:2675

msg.writeTag()



/*:199*/


//line amail.w:2820



/*152:*/


//line amail.w:2036

{
glog.V(debug).Infof("refresh a message '%d'\n",msg.id)
mrfch<-&refresh{seek,append(messages{},msg)}
}



/*:152*/


//line amail.w:2821

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


/*220:*/


//line amail.w:3151

{


/*225:*/


//line amail.w:3258



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3259




/*:225*/


//line amail.w:3153

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



/*:220*/


//line amail.w:2837

continue
case"Save":


/*230:*/


//line amail.w:3283

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




/*:230*/


//line amail.w:2840

continue
case"Q":
quote= true
if args:=strings.Fields(ev.Arg);len(args)> 0{
ev.Text= args[0]
ev.Arg= strings.Join(args," ")
}
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


/*231:*/


//line amail.w:3308

{


/*232:*/


//line amail.w:3348

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*233:*/


//line amail.w:3360

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


/*237:*/


//line amail.w:3451

{


/*240:*/


//line amail.w:3579



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3580




/*:240*/


//line amail.w:3453

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
case"to":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3478

to= append(to,f...)
case"cc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3481

cc= append(cc,f...)
case"bcc":


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3484

bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%s",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}


/*238:*/


//line amail.w:3565

for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
if sf:=strings.Fields(f[i]);len(sf)> 1{
f[i]= strings.TrimSpace(sf[len(sf)-1])
}
}




/*:238*/


//line amail.w:3499

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



/*:237*/


//line amail.w:3375

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:233*/


//line amail.w:3357




/*:232*/


//line amail.w:3310

name:=fmt.Sprintf("Amail/%s/%d/%sReply%s",
msg.box.name,
msg.id,
func()string{if quote{return"Q"};return""}(),
func()string{if replyall{return"all"};return""}())


/*52:*/


//line amail.w:789

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:52*/


//line amail.w:3316

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


/*234:*/


//line amail.w:3383

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


/*235:*/


//line amail.w:3401

{
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
buf= append(buf,'>',' ')
if strings.HasSuffix(s,"\r\n"){
s= strings.TrimRight(s,"\r\n")
s+= "\n"
}
buf= append(buf,s...)
if err==io.EOF{
break
}
}
}



/*:235*/


//line amail.w:3394

f.Close()
}else if len(msg.html)!=0{


/*236:*/


//line amail.w:3418

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


/*235:*/


//line amail.w:3401

{
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
buf= append(buf,'>',' ')
if strings.HasSuffix(s,"\r\n"){
s= strings.TrimRight(s,"\r\n")
s+= "\n"
}
buf= append(buf,s...)
if err==io.EOF{
break
}
}
}



/*:235*/


//line amail.w:3443

c1.Wait()
c2.Wait()
f.(io.Closer).Close()
}



/*:236*/


//line amail.w:3397

}



/*:234*/


//line amail.w:3337

}else{
buf= append(buf,fmt.Sprintf("Include: Mail/%s/%d/raw\n",msg.box.name,msg.id)...)

}
buf= append(buf,'\n')
w.Write(buf)


/*243:*/


//line amail.w:3594

if msg!=nil{
writeSignature(w,msg.box)
}else{
writeSignature(w,nil)
}



/*:243*/


//line amail.w:3344

}



/*:231*/


//line amail.w:2860

continue
case"Up":


/*131:*/


//line amail.w:1799

var parent*message
{
ch:=make(chan parentmsg)
glog.V(debug).Infof("getting parent for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
parent= <-ch
}



/*:131*/


//line amail.w:2863

if parent!=nil{


/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:2865

name:=parent.box.name
id:=parent.id


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:2868



/*88:*/


//line amail.w:1246

ach<-&struct{m msgmap;a action}{msgs,view}



/*:88*/


//line amail.w:2869

}
continue
case"Down":


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2873

if len(children)!=0{


/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:2875

name:=children[0].box.name
id:=children[0].id


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:2878



/*88:*/


//line amail.w:1246

ach<-&struct{m msgmap;a action}{msgs,view}



/*:88*/


//line amail.w:2879

}
continue
case"Prev":


/*200:*/


//line amail.w:2679

pmsg:=msg.prev()



/*:200*/


//line amail.w:2883

if pmsg!=nil{


/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:2885

name:=pmsg.box.name
id:=pmsg.id


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:2888



/*88:*/


//line amail.w:1246

ach<-&struct{m msgmap;a action}{msgs,view}



/*:88*/


//line amail.w:2889

}
continue
case"Next":


/*202:*/


//line amail.w:2702

nmsg:=msg.next()



/*:202*/


//line amail.w:2893

if nmsg!=nil{


/*85:*/


//line amail.w:1226

msgs:=make(msgmap)



/*:85*/


//line amail.w:2895

name:=nmsg.box.name
id:=nmsg.id


/*86:*/


//line amail.w:1230

glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n",id,name)
msgs[name]= append(msgs[name],id)



/*:86*/


//line amail.w:2898



/*88:*/


//line amail.w:1246

ach<-&struct{m msgmap;a action}{msgs,view}



/*:88*/


//line amail.w:2899

}
continue
}
}else if(ev.Type&goacme.Look)==goacme.Look{
}
msg.w.UnreadEvent(ev)

}
}()



/*:206*/


//line amail.w:2669

}
return
}



/*:198*/



/*201:*/


//line amail.w:2683

func(this*message)prev()(pmsg*message){
msg:=this


/*131:*/


//line amail.w:1799

var parent*message
{
ch:=make(chan parentmsg)
glog.V(debug).Infof("getting parent for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
parent= <-ch
}



/*:131*/


//line amail.w:2686

if parent==nil{
return
}
msg= parent


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2691

for _,v:=range children{
if v==this{
break
}
pmsg= v
}
return
}



/*:201*/



/*203:*/


//line amail.w:2706

func(this*message)next()(nmsg*message){
msg:=this


/*131:*/


//line amail.w:1799

var parent*message
{
ch:=make(chan parentmsg)
glog.V(debug).Infof("getting parent for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
parent= <-ch
}



/*:131*/


//line amail.w:2709

if parent==nil{
return
}
msg= parent


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2714

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



/*:203*/



/*204:*/


//line amail.w:2729

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
func()string{


/*131:*/


//line amail.w:1799

var parent*message
{
ch:=make(chan parentmsg)
glog.V(debug).Infof("getting parent for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
parent= <-ch
}



/*:131*/


//line amail.w:2745

if parent!=nil{
return"Up "
}
return""
}(),
func()string{


/*128:*/


//line amail.w:1757

var children idmessages
{
ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{msg*message;val interface{}}{msg,ch}
children= <-ch
}



/*:128*/


//line amail.w:2752

if len(children)!=0{
return"Down "
}
return""
}(),
func()string{


/*200:*/


//line amail.w:2679

pmsg:=msg.prev()



/*:200*/


//line amail.w:2759

if pmsg!=nil{
return"Prev "
}
return""
}(),
func()string{


/*202:*/


//line amail.w:2702

nmsg:=msg.next()



/*:202*/


//line amail.w:2766

if nmsg!=nil{
return"Next "
}
return""
}()));
err!=nil{
glog.Errorf("can't set a tag of the message window: %v",err)
}
}



/*:204*/



/*213:*/


//line amail.w:3002

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



/*:213*/



/*214:*/


//line amail.w:3052

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





/*:214*/



/*218:*/


//line amail.w:3102

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



/*:218*/



/*219:*/


//line amail.w:3121

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
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
if strings.HasSuffix(s,"\r\n"){
s= strings.TrimRight(s,"\r\n")
s+= "\n"
}
buf= append(buf,s...)
if err==io.EOF{
break
}
}
return buf,nil
}




/*:219*/



/*221:*/


//line amail.w:3184

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



/*:221*/



/*222:*/


//line amail.w:3202

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



/*:222*/



/*245:*/


//line amail.w:3608

func writeSignature(w*goacme.Window,box*mailbox){


/*216:*/


//line amail.w:3088



/*229:*/


//line amail.w:3279

once.Do(func(){

/*217:*/


//line amail.w:3092

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:217*/



/*226:*/


//line amail.w:3262

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:226*/



/*241:*/


//line amail.w:3583

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:241*/


//line amail.w:3280
})



/*:229*/


//line amail.w:3089




/*:216*/


//line amail.w:3610

var f io.ReadCloser
var err error
if box!=nil{
f,err= os.Open(fmt.Sprintf("%s/mail/%s.signature",home,box.name))
}
if err!=nil||f==nil{
f,err= os.Open(fmt.Sprintf("%s/mail/signature",home))
}
if err==nil{
w.Write([]byte("\n"))
b:=bufio.NewReader(f)
for buf,err:=b.ReadBytes('\n');err==nil||err==io.EOF;buf,err= b.ReadBytes('\n'){
w.Write(buf)
if err==io.EOF{
break
}
}
f.Close()
}


/*95:*/


//line amail.w:1317

glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:95*/


//line amail.w:3630

}



/*:245*/


