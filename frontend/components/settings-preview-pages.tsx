"use client";

import { Save, ServerCog, UserPlus } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

function PreviewNotice({ children }: { children: React.ReactNode }) {
  return (
    <div className="mb-5 flex items-center gap-2">
      <span className="rounded-full border border-primary/30 bg-primary/10 px-2.5 py-1 text-xs font-semibold text-primary">
        Preview only
      </span>
      <p className="text-xs text-muted-foreground">{children}</p>
    </div>
  );
}

export function GeneralSettingsPreview() {
  const [workspace, setWorkspace] = useState("RPMP Monitoring");
  const [timezone, setTimezone] = useState("asia-jakarta");
  const [retention, setRetention] = useState("90");
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [alerts, setAlerts] = useState(true);

  return (
    <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
      <PreviewNotice>
        Changes stay in this browser view and are not saved.
      </PreviewNotice>
      <div className="max-w-2xl space-y-6">
        <section className="surface-card space-y-5 p-6 hover:shadow-none">
          <h2 className="text-sm font-bold">Workspace</h2>
          <div className="space-y-2">
            <Label htmlFor="workspace-name">Workspace name</Label>
            <Input
              id="workspace-name"
              value={workspace}
              onChange={(event) => setWorkspace(event.target.value)}
            />
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="workspace-timezone">Timezone</Label>
              <Select
                value={timezone}
                onValueChange={(value) => value && setTimezone(value)}
              >
                <SelectTrigger id="workspace-timezone" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="asia-jakarta">
                    Asia/Jakarta (GMT+7)
                  </SelectItem>
                  <SelectItem value="asia-singapore">
                    Asia/Singapore (GMT+8)
                  </SelectItem>
                  <SelectItem value="utc">UTC</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="data-retention">Data retention</Label>
              <Select
                value={retention}
                onValueChange={(value) => value && setRetention(value)}
              >
                <SelectTrigger id="data-retention" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="30">30 days</SelectItem>
                  <SelectItem value="90">90 days</SelectItem>
                  <SelectItem value="365">1 year</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </section>

        <section className="surface-card space-y-1 p-6 hover:shadow-none">
          <h2 className="mb-4 text-sm font-bold">Monitoring</h2>
          <label className="flex cursor-pointer items-center justify-between gap-4 border-b border-border py-3">
            <span>
              <span className="block text-sm font-medium">
                Live auto-refresh
              </span>
              <span className="block text-xs text-muted-foreground">
                Preview a dashboard refresh preference.
              </span>
            </span>
            <Switch
              aria-label="Live auto-refresh"
              checked={autoRefresh}
              onCheckedChange={setAutoRefresh}
            />
          </label>
          <label className="flex cursor-pointer items-center justify-between gap-4 py-3">
            <span>
              <span className="block text-sm font-medium">Failure alerts</span>
              <span className="block text-xs text-muted-foreground">
                Preview an alert preference for unhealthy use cases.
              </span>
            </span>
            <Switch
              aria-label="Failure alerts"
              checked={alerts}
              onCheckedChange={setAlerts}
            />
          </label>
        </section>

        <div className="flex justify-end">
          <Button
            className="rounded-full"
            onClick={() =>
              toast.success("Settings save simulated", {
                description:
                  "Your changes remain in this preview only. Nothing was saved.",
              })
            }
          >
            <Save className="mr-2 h-4 w-4" />
            Simulate save
          </Button>
        </div>
      </div>
    </main>
  );
}

type Member = {
  id: string;
  name: string;
  email: string;
  role: "Admin" | "Operator" | "Viewer";
  status: "Active" | "Invited";
};

const initialMembers: Member[] = [
  {
    id: "u-1",
    name: "RPMP Admin",
    email: "admin@example.com",
    role: "Admin",
    status: "Active",
  },
  {
    id: "u-2",
    name: "Dewi Lestari",
    email: "dewi.lestari@example.com",
    role: "Operator",
    status: "Active",
  },
  {
    id: "u-3",
    name: "Rizky Pratama",
    email: "rizky.p@example.com",
    role: "Viewer",
    status: "Active",
  },
  {
    id: "u-4",
    name: "Sinta Maharani",
    email: "sinta.m@example.com",
    role: "Operator",
    status: "Invited",
  },
];

export function UserManagementPreview() {
  const [members, setMembers] = useState(initialMembers);
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<Member["role"]>("Viewer");

  function invite() {
    if (!/^\S+@\S+\.\S+$/.test(email)) {
      toast.error("Enter a valid email address.");
      return;
    }

    setMembers((current) => [
      ...current,
      {
        id: `u-${current.length + 1}`,
        name: email.split("@")[0] ?? email,
        email,
        role,
        status: "Invited",
      },
    ]);
    toast.success("Invitation simulated", {
      description: `${email} was added to this preview as ${role}. No invitation was sent and no user was created.`,
    });
    setEmail("");
  }

  return (
    <main className="flex-1 space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <PreviewNotice>
        Members are sample data. Role changes and invitations are not saved or
        sent.
      </PreviewNotice>
      <section className="surface-card p-6 hover:shadow-none">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/12 text-primary">
            <UserPlus className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-sm font-bold">Preview an invitation</h2>
            <p className="text-xs text-muted-foreground">
              Add a sample member to this browser view. No account or email is
              created.
            </p>
          </div>
        </div>
        <div className="mt-5 grid gap-3 sm:grid-cols-[minmax(0,1fr)_160px_auto]">
          <Input
            aria-label="Preview member email"
            placeholder="name@example.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
          />
          <Select
            value={role}
            onValueChange={(value) => setRole(value as Member["role"])}
          >
            <SelectTrigger aria-label="Preview member role">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Admin">Admin</SelectItem>
              <SelectItem value="Operator">Operator</SelectItem>
              <SelectItem value="Viewer">Viewer</SelectItem>
            </SelectContent>
          </Select>
          <Button className="rounded-full" onClick={invite}>
            Simulate invite
          </Button>
        </div>
      </section>

      <section className="surface-card overflow-hidden hover:shadow-none">
        <div className="border-b border-border px-6 py-4">
          <h2 className="text-sm font-bold">
            Sample members ({members.length})
          </h2>
        </div>
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="min-w-56">User</TableHead>
              <TableHead>Preview role</TableHead>
              <TableHead>Sample status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {members.map((member) => (
              <TableRow key={member.id}>
                <TableCell>
                  <div className="flex items-center gap-3">
                    <Avatar className="h-8 w-8">
                      <AvatarFallback className="bg-primary/15 text-xs font-bold text-primary">
                        {member.name.slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0">
                      <p className="truncate text-sm font-semibold">
                        {member.name}
                      </p>
                      <p className="truncate text-xs text-muted-foreground">
                        {member.email}
                      </p>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <Select
                    value={member.role}
                    onValueChange={(value) =>
                      setMembers((current) =>
                        current.map((item) =>
                          item.id === member.id
                            ? { ...item, role: value as Member["role"] }
                            : item,
                        ),
                      )
                    }
                  >
                    <SelectTrigger
                      className="h-8 w-32 text-xs"
                      aria-label={`Preview role for ${member.name}`}
                    >
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Admin">Admin</SelectItem>
                      <SelectItem value="Operator">Operator</SelectItem>
                      <SelectItem value="Viewer">Viewer</SelectItem>
                    </SelectContent>
                  </Select>
                </TableCell>
                <TableCell>
                  <span
                    className={
                      member.status === "Active"
                        ? "rounded-full bg-success/12 px-2.5 py-0.5 text-[11px] font-semibold text-success"
                        : "rounded-full bg-warning/15 px-2.5 py-0.5 text-[11px] font-semibold text-warning"
                    }
                  >
                    {member.status}
                  </span>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </section>
    </main>
  );
}

export function EmailSettingsPreview() {
  const [host, setHost] = useState("smtp.example.com");
  const [port, setPort] = useState("587");
  const [encryption, setEncryption] = useState("tls");
  const [senderName, setSenderName] = useState("RPMP Reports");
  const [senderEmail, setSenderEmail] = useState("no-reply@example.com");
  const [enabled, setEnabled] = useState(true);

  return (
    <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
      <PreviewNotice>
        This sample email configuration stays in local state. RPMP will not
        connect to a server or send email.
      </PreviewNotice>
      <div className="max-w-2xl space-y-6">
        <section className="surface-card space-y-5 p-6 hover:shadow-none">
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/12 text-primary">
                <ServerCog className="h-5 w-5" />
              </div>
              <div>
                <h2 className="text-sm font-bold">
                  Sample email server
                </h2>
                <p className="text-xs text-muted-foreground">
                  Display values only. No connection is made.
                </p>
              </div>
            </div>
            <label className="flex items-center gap-2 text-xs text-muted-foreground">
              <Switch
                aria-label="Email configuration enabled"
                checked={enabled}
                onCheckedChange={setEnabled}
              />
              {enabled ? "Enabled in preview" : "Disabled in preview"}
            </label>
          </div>

          <div className="grid gap-4 sm:grid-cols-[minmax(0,1fr)_120px_140px]">
            <div className="space-y-2">
              <Label htmlFor="email-host">Host</Label>
              <Input
                id="email-host"
                value={host}
                onChange={(event) => setHost(event.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email-port">Port</Label>
              <Input
                id="email-port"
                value={port}
                onChange={(event) => setPort(event.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email-encryption">Encryption</Label>
              <Select
                value={encryption}
                onValueChange={(value) => value && setEncryption(value)}
              >
                <SelectTrigger id="email-encryption" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="tls">TLS</SelectItem>
                  <SelectItem value="ssl">SSL</SelectItem>
                  <SelectItem value="none">None</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="sender-name">Sender name</Label>
              <Input
                id="sender-name"
                value={senderName}
                onChange={(event) => setSenderName(event.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="sender-email">Sender email</Label>
              <Input
                id="sender-email"
                value={senderEmail}
                onChange={(event) => setSenderEmail(event.target.value)}
              />
            </div>
          </div>
        </section>

        <div className="flex flex-col justify-end gap-2 sm:flex-row">
          <Button
            variant="outline"
            className="rounded-full"
            onClick={() =>
              toast.success("Email test simulated", {
                description: `No test email was sent from ${senderEmail}.`,
              })
            }
          >
            Simulate test email
          </Button>
          <Button
            className="rounded-full"
            onClick={() =>
              toast.success("Email settings save simulated", {
                description: `${host}:${port} (${encryption.toUpperCase()}) remains in this preview only. Nothing was saved.`,
              })
            }
          >
            <Save className="mr-2 h-4 w-4" />
            Simulate save
          </Button>
        </div>
      </div>
    </main>
  );
}
