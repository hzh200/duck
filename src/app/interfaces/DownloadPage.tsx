import React, { useEffect, useState } from "react";
import { Extractor, extractors } from "@/lib/extractors";

import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

import { Input } from "@/components/ui/input";

function DownloadPage(setting: {[key: string]: any}) {
  const [choosenExtractor, setChoosenExtractor] = useState<Extractor>(extractors[0]);
  const [url, setUrl] = useState<string>("");
  const [options, setOptions] = useState<{[key: string]: any}>({});
  const [optionTypes, setOptionTypes] = useState<{[key: string]: any}>({});
  const [location, setLocation] = useState<string>(setting["downloadDirectory"]);
  const [errMessage, setErrMessage] = useState<string>("");
  const [extracting, setExtracting] = useState<boolean>(false);
  const [downloading, setDownloading] = useState<boolean>(false);
  const [extracted, setExtracted] = useState<boolean>(false);

  const extract = () => {
    if (url === "") {
        return;
    }

    setDownloading(false);
    setExtracting(true);
    setExtracted(false);
    setErrMessage("");
    fetch(`http://127.0.0.1:${setting["kernelPort"]}/extract?url=${url}&extractor=${choosenExtractor.name}`).then(res => res.json()).then(async res => {
        if (res["errMessage"]) {
            throw new Error(res["errMessage"]);
        }
        
        setOptions(res["options"]);
        setOptionTypes(res["optionTypes"]);
        setExtracting(false);
        setExtracted(true);
        setErrMessage("");
    }).catch((err: Error) => {
        setOptions({});
        setOptionTypes({});
        setExtracting(false);
        setExtracted(false);
        setErrMessage(err.name + ":" + err.message + "\n" + err.stack);
    });
  };

  const download = () => {
    setDownloading(true)
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    fetch('http://127.0.0.1:9000/download', {
        method: "POST",
        body: JSON.stringify({
            ...options,
            ...{
                "Url": url,
                "Extractor": choosenExtractor.name,
                "Location": location
            }
        }),
        headers: myHeaders,
    }).then(res => res.json()).then(res => {
        if (res["status"] === "succeed") {
            setDownloading(false);
            toast("Adding task/taskSet succeed", {
                description: JSON.stringify(options),
                action: {
                  label: "Got it",
                  onClick: () => console.log("Got it"),
                }
            });
        } else {
            setErrMessage(res["errMessage"]);
        }
    });
  };

  return (
    <div className='h-full p-2 flex items-center justify-center'>
        <Card className="w-10/12 m-2">
            <CardHeader>
                <CardTitle>Download</CardTitle>
                <CardDescription>Choose an extractor and type in the URL.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="grid w-full items-center gap-4 space-y-1.5">
                    <div className="flex flex-col w-full space-y-1.5">
                        <Label htmlFor="extractor">Extractor</Label>
                        <Select onValueChange={(value) => {
                            for (const extractor of extractors) {
                                if (extractor.name === value) {
                                    setChoosenExtractor(extractor)
                                }
                            }
                        }} defaultValue={choosenExtractor.name}>
                            <SelectTrigger id="extractor" className="w-[120px]">
                                <SelectValue placeholder="Select a extractor" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectLabel>Extractor</SelectLabel>
                                    {extractors.map((extractor: Extractor, index) => (
                                        <SelectItem key={index} value={extractor.name}>{extractor.name}</SelectItem>
                                    ))}
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                        <Label htmlFor="url">URL</Label>
                        <Textarea id="url" className="h-[60px] w-full" placeholder="Type url here." value={url} onChange={event => setUrl(event.target.value)} />
                        <div className="flex flex-row-reverse justify-between">
                            <Button className={`h-10 w-20 text-base`} onClick={() => extract()}>
                                Extract
                            </Button>
                            <div>
                                {extracting ? <span><div className="loading" />extracting</span> : ""}
                            </div> 
                        </div>                      
                        <div className={`flex flex-col space-y-1.5 ${extracted && errMessage === "" ? "" : "hidden"}`}>
                            <Label htmlFor="location">Location</Label>
                            <Input id="location" value={location} onChange={event => setLocation(event.target.value)} />
                            {Object.entries(options).map(([key, value], i) => {
                                if (optionTypes[key] === "Input") {
                                    return (
                                        <React.Fragment key={i}>
                                            <Label htmlFor={key}>{key}</Label>
                                            <Input id={key} value={options[key]} onChange={event => {
                                                const newOptions = {...options};
                                                const newVal = event.target.value;
                                                newOptions[key] = newVal;
                                                setOptions(newOptions);
                                            }} />
                                        </React.Fragment>
                                    )
                                } else if (optionTypes[key] === "") {

                                }
                            })}
                        </div>
                        <div className="flex flex-row-reverse justify-between">
                            <Button className={`h-10 w-20 text-base ${extracted ? "" : "hidden"}`} onClick={() => download()}>
                                Download
                            </Button>
                            <span>
                                {downloading ? <span><div className="loading" />downloading</span> : ""}
                            </span> 
                        </div>
                        <Card className={`h-[150px] w-full p-3 ${errMessage !== "" ? "" : "hidden"}`}>
                            {errMessage}
                        </Card>  
                    </div>
                </div>
            </CardContent>
            <CardFooter className="">
            </CardFooter>
        </Card>

        
    </div>
  );
}

export default DownloadPage;
