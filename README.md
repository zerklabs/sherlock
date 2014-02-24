sherlock
========

__this code is a proof-of-concept__

Sherlock helps rank "important" words in a given corpus. For instance, when given the output of `strings -td -eS malware.exe`, you get something similar to:

(see [malware_strings.txt][malware] for the longer version)
```
[...]
146460 application/x-www-form-urlencoded
146500 Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)
146552 windows/cartoon
146568 POST
146576 compose.aspx?s=%4X%4X%4X%4X%4X%4X
146612 %s ip port [tcp|http|udp]
146716 www.google.com
146736 http
146867 qs*¹>
146894 0@íµ ÷Æ°>
146908 €„.A
146940 \wship6
146948 \ws2_32
146956 freeaddrinfo
146972 getnameinfo
146984 getaddrinfo
147004 65535
147152 ÿÿÿÿ 9A
147234 ÀÿÿÿßA%s:%s:%lld
147392 ÿÿÿÿ
147408 ÿÿÿÿÐ
147436 ÿÿÿÿ
151568 ÿÿÿÿ`
[...]
152000 ÿÿÿÿÐ
152120 ÿÿÿÿ
152128 ÿÿÿÿ
152720 ÿÿÿÿ
[...]
153016 ÿÿÿÿ
153032 ÿÿÿÿ
153112 ÿÿÿÿ
153128 ÿÿÿÿ
[...]
153664 ÿÿÿÿ°"B
153720 ÿÿÿÿà"B
153760 ÿÿÿÿ
153800 ÿÿÿÿ #B
[...]
155006 WSASendTo
155018 WSARecvFrom
155030 WS2_32.dll
155044 FreeLibrary
155058 GetProcAddress
155076 LoadLibraryA
155092 CloseHandle
155106 WriteFile
155118 CreateFileA
155132 GetTempFileNameA
155152 GetTempPathA
155168 WaitForMultipleObjects
155194 DeleteFileW
155208 Sleep
[...]
```

 and then once we score and classify each word, we get the output:

```
[...]
1,5.2,application/x-www-form-urlencoded
1,0.79999995,Mozilla/4.0
1,1.6000001,(compatible;
1,1.2,Windows
1,2.5000002,windows/cartoon
1,1.9000002,compose.aspx?s=%4X%4X%4X%4X%4X%4X
1,1.4000001,[tcp|http|udp]
1,2.0000002,www.google.com
1,2.2000003,freeaddrinfo
1,2.0000002,getnameinfo
1,2.0000002,getaddrinfo
1,1.6000001,WSASendTo
1,2.0000002,WSARecvFrom
1,2.0000002,FreeLibrary
1,2.6000004,GetProcAddress
1,2.2000003,LoadLibraryA
1,2.0000002,CloseHandle
1,1.6000001,WriteFile
1,2.0000002,CreateFileA
1,3.0000005,GetTempFileNameA
1,2.2000003,GetTempPathA
1,4.2000003,WaitForMultipleObjects
1,2.0000002,DeleteFileW
1,0.8,Sleep
[...]
```

where the format of the output is:

CLASSIFICATION,SCORE,WORD

CLASSIFICATION can be:
  * -1 = garbage
  * 0 = neutral
  * 1 = important


[malware]: https://github.com/zerklabs/sherlock/tree/master/docs
