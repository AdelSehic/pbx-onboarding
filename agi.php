#!/usr/bin/php -q
<?php
require('phpagi.php');

$agi = new AGI();
$caller = $agi->get_variable('CALLERID(num)')['data'];
$pass = $agi->database_get("password",$caller);

$agi->answer();
$agi->exec("PLAYBACK", "welcome");

pwcheck($pass, $agi);

ivr:
$option = $agi->get_data("sekretarica", 20000, 1);
switch ($option['result']){
    case 1:
        $exten = $agi->get_data("vm-extension", 20000);
        if ($exten['result']==$caller){
            $agi->exec("PLAYBACK", "invalid");
            $agi->hangup();
        }
        $agi->set_variable("tocall", $exten['result']);
        $agi->set_variable("redirect", "1");
        break;
    case 2:
        pwcheck($pass, $agi);
        $pw1 = $agi->get_data("vm-newpassword", 20000, 4);
        $pw2 = $agi->get_data("vm-reenterpassword", 20000, 4);
        if ($pw1 != $pw2 ){
            $agi->exec("PLAYBACK", "vm-mismatch");
            $agi->hangup();
        }
        $agi->database_put("password", $caller, $pw1['result']);
        $agi->exec("PLAYBACK", "vm-passchanged");
        break;
    case 3:
        goto ivr;
}

function pwcheck($pass, $agi){
    if(!$pass['result']) return;

    $result = $agi->get_data("vm-password", 9000, 4);
    if ($result['result'] != $pass['data']) {
        $agi->exec("PLAYBACK", "auth-incorrect");
        $agi->hangup();
    }
}
?>



/* DIALPLAN

    [sets]

    exten => 101,1,Dial(PJSIP/SOFTPHONE_A)
    exten => 102,1,Dial(PJSIP/SOFTPHONE_B)

    exten => 123,1,Noop(php ivr channel)
        same => n,AGI(agi.php)
        same => n,GotoIf($["${redirect}" = "1"]?${tocall},1:)
        same => n,Hangup()

    exten => i,1,Playback(invalid)

*/