#!/usr/bin/env perl

use strict;
use warnings;

binmode STDOUT, ':utf8';

use Pod::Markdown;

my ($file) = @ARGV;

my $pod_string = do { local $/; open my $fh, '<', $file or die $!; <$fh> };

$pod_string =~ s{^(.*?)\r?\n\r?\n}{}ms;
my $meta = $1;

$pod_string =~ s{^\[cut\]}{##########}ms;

$pod_string = "=pod\n\n" . $pod_string;

my $markdown;
my $parser = Pod::Markdown->new;
$parser->output_string(\$markdown);
$parser->parse_string_document($pod_string);

$markdown =~ s{\\##########}{[cut]};

print $meta, "\n\n";
print $markdown;
