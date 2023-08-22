#!/usr/bin/php -q
<?php

$weatherURL="http://www.nws.noaa.gov/data/current_obs/KMDQ.xml";

set_time_limit(60);
ob_implicit_flush(false);

if (!defined('STDIN')){
    define('STDIN', fopen('php://stdin', 'r'));
}
if (!defined('STDOUT')){
    define('STDOUT', fopen('php://stdout', 'w'));
}
if (!defined('STDERR')){
    define('STDERR', fopen('php://stderr', 'w'));
}

while (!feof(STDIN)){
    $temp = trim(fgets(STDIN,4096));
    if (($temp == '') || ($temp == '\n')){
        break;
    }
    $s = explode(":", $temp);
    $name = str_replace("agi_","",$s[0]);
    $agi[$name] = trim($s[1]);
}

foreach($agi as $key=>$value)
{
    fwrite(STDERR,"-- $key = $value\n");
    fflush(STDERR);
}

$weatherPage=file_get_contents($weatherURL);

if (preg_match("/<temp_f>([0-9]+)<\/temp_f>/i",$weatherPage,$matches))
{
    $currentTemp=$matches[1];
}

if (preg_match("/<wind_dir>North<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='northerly';
}
elseif (preg_match("/<wind_dir>South<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='southerly';
}
elseif (preg_match("/<wind_dir>East<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='easterly';
}
elseif (preg_match("/<wind_dir>West<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='westerly';
}
elseif (preg_match("/<wind_dir>Northwest<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='northwesterly';
}
elseif (preg_match("/<wind_dir>Northeast<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='northeasterly';
}
elseif (preg_match("/<wind_dir>Southwest<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='southwesterly';
}
elseif (preg_match("/<wind_dir>Southeast<\/wind_dir>/i",$weatherPage))
{
    $currentWindDirection='southeasterly';
}

if (preg_match("/<wind_mph>([0-9.]+)<\/wind_mph>/i",$weatherPage,$matches))
{
    $currentWindSpeed = $matches[1];
}

if ($currentTemp)
{
    fwrite(STDOUT,"STREAM FILE temperature \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"STREAM FILE is \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"SAY NUMBER $currentTemp \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"STREAM FILE degrees \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"STREAM FILE fahrenheit \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
}

if ($currentWindDirection && $currentWindSpeed)
{
    fwrite(STDOUT,"STREAM FILE with \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"STREAM FILE $currentWindDirection \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"STREAM FILE wx/winds \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"STREAM FILE at \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite(STDOUT,"SAY NUMBER $currentWindSpeed \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
    fwrite($STDOUT,"STREAM FILE miles-per-hour \"\"\n");
    fflush(STDOUT);
    $result = trim(fgets(STDIN,4096));
    checkresult($result);
}

function checkresult($res)
{
    trim($res);
    if (preg_match('/^200/',$res))
    {
        if (! preg_match('/result=(-?\d+)/',$res,$matches))
        {
            fwrite(STDERR,"FAIL ($res)\n");
            fflush(STDERR);
            return 0;
        }
        else
        {
            fwrite(STDERR,"PASS (".$matches[1].")\n");
            fflush(STDERR);
            return $matches[1];
        }
    }
    else
    {
        fwrite(STDERR,"FAIL (unexpected result '$res')\n");
        fflush(STDERR);
        return -1;
    }
}
?>