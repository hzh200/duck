import React, { useState } from "react";
import { Parser, parsers } from "@/lib/parsers";

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

function DownloadPage() {
  const [choosenParser, setChoosenParser] = useState<Parser>(parsers[0]);
  const [parsed, setParsed] = useState<boolean>(false);
  const [options, setOptions] = useState<Array<[string, any]>>([]);

  const download = () => {
    toast("Event has been created", {
        description: "Sunday, December 03, 2023 at 9:00 AM",
        action: {
          label: "Undo",
          onClick: () => console.log("Undo"),
        },
    })
  };

  return (
    <div className='h-full p-2 flex items-center justify-center'>
        <Card className="w-fit m-2">
            <CardHeader>
                <CardTitle>Parse</CardTitle>
                <CardDescription>Choose a parser and type in the URL.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="grid w-full items-center gap-4">
                    <div className="flex flex-col space-y-1.5">
                        <Label htmlFor="parser">Parser</Label>
                        <Select onValueChange={(value) => {
                            for (const parser of parsers) {
                                if (parser.name === value) {
                                    setChoosenParser(parser)
                                }
                            }
                        }} defaultValue={choosenParser.name}>
                            <SelectTrigger id="parser" className="w-[120px]">
                                <SelectValue placeholder="Select a parser" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                <SelectLabel>Parser</SelectLabel>
                                {parsers.map((parser: Parser, index) => (
                                    <SelectItem key={index} value={parser.name}>{parser.name}</SelectItem>
                                ))}
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                    <div className="flex flex-col space-y-1.5 w-[650px]">
                        <Label htmlFor="url">URL</Label>
                        <Textarea id="url" className="h-[80px]" placeholder="Type url here." />
                    </div>
                </div>
            </CardContent>
            <CardFooter className="flex justify-between">
                <div>{/* <Button variant="outline">Reset</Button> */}</div>
                <Button className='h-12 w-20 text-base' onClick={() => setParsed(true)}>
                    Parse
                </Button>
            </CardFooter>
        </Card>
        <Card className={`w-full m-2 ${parsed ? "" : "hidden"}`}>
            <CardHeader>
                <CardTitle>Download</CardTitle>
                <CardDescription>Choose wanted formats.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="grid w-full items-center gap-4">
                    { options.map(item => (
                        <div className="flex flex-col space-y-1.5">
                            <Label htmlFor={item[0]}>{item[0]}</Label>
                            <Textarea id={item[0]} placeholder="Input here." />
                        </div>
                    )) }
                </div>
            </CardContent>
            <CardFooter className="flex justify-between">
                <div>{/* <Button variant="outline">Reset</Button> */}</div>
                <Button className='h-12 w-22 text-base' onClick={() => download()}>
                    Download
                </Button>
            </CardFooter>
        </Card>
    </div>
  );
}

export default DownloadPage;
